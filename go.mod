module github.com/lbryio/transcoder

go 1.16

require (
	github.com/alecthomas/kong v0.2.12
	github.com/aws/aws-sdk-go v1.36.29
	github.com/brk0v/directio v0.0.0-20190225130936-69406e757cf7
	github.com/c2h5oh/datasize v0.0.0-20200825124411-48ed595a09d2
	github.com/draganm/miniotest v0.1.0
	github.com/fasthttp/router v1.3.3
	github.com/floostack/transcoder v1.2.0
	github.com/grafov/m3u8 v0.11.1
	github.com/karlseguin/ccache/v2 v2.0.8
	github.com/karrick/godirwalk v1.16.1
	github.com/lbryio/lbry.go/v2 v2.7.2-0.20210416195322-6516df1418e3
	github.com/mattn/go-sqlite3 v1.14.4
	github.com/nikooo777/lbry-blobs-downloader v1.0.4
	github.com/orcaman/concurrent-map v0.0.0-20190826125027-8c72a8bb44f6
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.5.0
	github.com/prometheus/client_golang v1.9.0
	github.com/simukti/sqldb-logger v0.0.0-20201125162808-c35f87e285f2
	github.com/simukti/sqldb-logger/logadapter/zapadapter v0.0.0-20201125162808-c35f87e285f2
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/valyala/fasthttp v1.17.0
	go.uber.org/goleak v1.1.10
	go.uber.org/zap v1.16.0
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/tools v0.1.0 // indirect
)

// replace github.com/floostack/transcoder => github.com/andybeletsky/transcoder v1.2.1

replace github.com/nikooo777/lbry-blobs-downloader => github.com/andybeletsky/lbry-blobs-downloader v1.0.4-fixed2

replace github.com/floostack/transcoder => /Users/silence/Documents/LBRY/ExternalForks/transcoder
