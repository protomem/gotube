package blobstore

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/protomem/gotube/internal/config"
)

type Storage struct {
	*minio.Client

	conf config.Blob
}

func New(conf config.Config) *Storage {
	return &Storage{conf: conf.Blob}
}

func (s *Storage) Connect(_ context.Context) error {
	var err error

	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(s.conf.AccessKey, s.conf.SecretKey, ""),
		Secure: s.conf.Secure,
	}

	s.Client, err = minio.New(s.conf.Addr, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Disconnect(_ context.Context) error {
	return nil
}
