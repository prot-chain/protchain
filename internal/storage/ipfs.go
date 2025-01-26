package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"protchain/internal/config"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/multiformats/go-multiaddr"
)

type IPFSStorage struct {
	client *rpc.HttpApi
}

func NewIPFSStorage(cfg *config.Config) *IPFSStorage {
	ma, err := multiaddr.NewMultiaddr(cfg.IPFSAddress)
	if err != nil {
		log.Fatalf("invalid multiaddress format for IPFS node: %w", err)
		return nil
	}

	client, err := rpc.NewApi(ma)
	if err != nil {
		log.Fatalf("unable to connect to IPFS node at %s: %w", cfg.IPFSAddress, err)
		return nil
	}
	return &IPFSStorage{client: client}
}

func (s *IPFSStorage) Upload(fileName string, fileData []byte) (string, error) {
	node := files.NewReaderFile(bytes.NewReader(fileData))
	cid, err := s.client.Unixfs().Add(context.Background(), node)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to IPFS: %w", err)
	}
	return cid.String(), nil
}

func (s *IPFSStorage) Download(fileName string) (io.ReadCloser, error) {
	filePath, err := path.NewPath(fileName)
	if err != nil {
		return nil, err
	}
	fmt.Println(filePath)
	node, err := s.client.Unixfs().Get(context.Background(), filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from IPFS: %w", err)
	}

	// Convert the files.Node to an io.ReadCloser
	reader := node.(files.File)
	if reader == nil {
		return nil, fmt.Errorf("failed to convert IPFS node to reader")
	}
	return io.NopCloser(reader), nil
}
