package storage

import "io"

type StorageI interface {
	Upload(fileName string, fileData []byte) (string, error)
	Download(fileName string) (io.ReadCloser, error)
}
