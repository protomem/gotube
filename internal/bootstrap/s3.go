package bootstrap

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Options struct {
	Addr   string
	Access string
	Secret string
}

func S3(ctx context.Context, opts S3Options) (*minio.Client, error) {
	const op = "bootstrap.S3"

	client, err := minio.New(opts.Addr, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.Access, opts.Secret, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return client, nil
}
