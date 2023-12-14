package bootstrap

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Options struct {
	Addr   string
	Key    string
	Secret string
	Secure bool
}

func S3(_ context.Context, opts S3Options) (*minio.Client, error) {
	const op = "bootstrap.S3"

	mopts := &minio.Options{
		Creds:  credentials.NewStaticV4(opts.Key, opts.Secret, ""),
		Secure: opts.Secure,
	}

	client, err := minio.New(opts.Addr, mopts)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", op, err)
	}

	return client, nil
}
