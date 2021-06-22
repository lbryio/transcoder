package manager

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/nikooo777/lbry-blobs-downloader/downloader"
	"github.com/nikooo777/lbry-blobs-downloader/shared"

	ljsonrpc "github.com/lbryio/lbry.go/v2/extras/jsonrpc"
	"github.com/lbryio/transcoder/pkg/mfr"
	"github.com/lbryio/transcoder/pkg/timer"
)

var (
	lbrytvAPI  = "https://api.lbry.tv/api/v1/proxy"
	cdnServer  = "https://cdn.lbryplayer.xyz/api/v3/streams"
	blobServer = "cdn.lbryplayer.xyz"
	udpPort    = 5568
	tcpPort    = 5567

	lbrytvClient   = ljsonrpc.NewClient(lbrytvAPI)
	downloadClient = &http.Client{
		Timeout: 1200 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 120 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   30 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
		},
	}
)

type WriteCounter struct {
	Loaded, Size   uint64
	Started        time.Time
	URL            string
	progressLogged map[int]bool
}

type TranscodingRequest struct {
	queue *mfr.Queue

	URI, Name, ClaimID, SDHash, ChannelURI, NormalizedName string
	ChannelSupportAmount                                   int64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Loaded += uint64(n)
	progress := int(float64(wc.Loaded) / float64(wc.Size) * 100)

	if progress > 0 && progress%25 == 0 && !wc.progressLogged[progress] {
		wc.progressLogged[progress] = true
		rate := int64(float64(wc.Loaded) / time.Since(wc.Started).Seconds())
		logger.Debugw(
			"download progress",
			"url", wc.URL,
			"size", wc.Size, "progress", int(progress), "rate", rate)
	}
	return n, nil
}

func ResolveRequest(uri string) (*TranscodingRequest, error) {
	c, err := Resolve(uri)
	if err != nil {
		return nil, err
	}

	if c.SigningChannel == nil {
		return nil, ErrNoSigningChannel
	}

	src := c.Value.GetStream().GetSource()
	if src == nil {
		return nil, errors.New("stream doesn't have source data")
	}
	h := hex.EncodeToString(src.SdHash)

	ch := strings.Replace(strings.ToLower(c.SigningChannel.CanonicalURL), "#", ":", 1)
	sup, _ := strconv.ParseFloat(c.SigningChannel.Meta.SupportAmount, 64)

	r := &TranscodingRequest{
		URI:                  uri,
		SDHash:               h,
		Name:                 c.Name,
		NormalizedName:       c.NormalizedName,
		ClaimID:              c.ClaimID,
		ChannelURI:           ch,
		ChannelSupportAmount: int64(math.Floor(sup)),
	}
	return r, nil
}

func Resolve(uri string) (*ljsonrpc.Claim, error) {
	resolved, err := lbrytvClient.Resolve(uri)
	if err != nil {
		return nil, err
	}

	c, ok := (*resolved)[uri]
	if !ok || c.CanonicalURL == "" {
		return nil, ErrStreamNotFound
	}
	return &c, nil
}

// Download retrieves a video stream from the lbrytv CDN and saves it to a temporary file.
func (c *TranscodingRequest) Download(dest string) (*os.File, int64, error) {
	shared.ReflectorPeerServer = fmt.Sprintf("%s:%d", blobServer, tcpPort)
	shared.ReflectorQuicServer = fmt.Sprintf("%s:%d", blobServer, udpPort)

	var readLen int64
	logger.Infow("downloading stream", "url", c.URI)
	t := timer.Start()

	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return nil, 0, err
	}

	if err := downloader.DownloadAndBuild(c.SDHash, false, downloader.UDP, c.streamFileName(), dest); err != nil {
		return nil, 0, err
	}
	t.Stop()

	fi, err := os.Stat(path.Join(dest, c.streamFileName()))
	if err != nil {
		return nil, 0, err
	}
	readLen = fi.Size()

	f, err := os.Open(path.Join(dest, c.streamFileName()))
	if err != nil {
		return nil, 0, err
	}

	rate := int64(float64(readLen) / t.Duration())
	logger.Infow("stream downloaded", "url", c.URI, "rate", rate, "size", readLen, "seconds_spent", t.DurationInt())
	return f, readLen, nil
}

func (r *TranscodingRequest) Release() {
	logger.Infow("transcoding request released", "lbry_url", r.URI)
	r.queue.Release(r.URI)
}

func (r *TranscodingRequest) Reject() {
	logger.Infow("transcoding request rejected", "lbry_url", r.URI)
	r.queue.Done(r.URI)
}

func (r *TranscodingRequest) Complete() {
	logger.Infow("transcoding request completed", "lbry_url", r.URI)
	r.queue.Done(r.URI)
}

func (c *TranscodingRequest) streamFileName() string {
	return c.SDHash
}

func SetCDNServer(s string) {
	cdnServer = s
}
