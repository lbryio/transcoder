package storage

import (
	"os"
	"path"
)

type LocalStorage struct {
	path string
}

func Local(path string) LocalStorage {
	return LocalStorage{path}
}

func (s LocalStorage) New(sdHash string) *LocalStream {
	return &LocalStream{rootPath: s.path, sdHash: sdHash}
}

func (s LocalStorage) Open(sdHash string) (*LocalStream, error) {
	ls := &LocalStream{rootPath: s.path, sdHash: sdHash}
	_, err := os.Stat(ls.FullPath())
	return ls, err
}

func (s LocalStorage) Delete(sdHash string) error {
	return os.RemoveAll(path.Join(s.path, sdHash))
}
