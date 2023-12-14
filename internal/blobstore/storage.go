package blobstore

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	*minio.Client
}

func New(addr, accessKey, secretKey string, secure bool) (*Storage, error) {
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	}

	client, err := minio.New(addr, opts)
	if err != nil {
		return nil, err
	}

	return &Storage{client}, nil
}
