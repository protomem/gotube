package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/protomem/gotube/internal/bootstrap"
	"github.com/protomem/gotube/pkg/logging"
)

var _ Storage = (*S3)(nil)

type S3 struct {
	logger logging.Logger
	client *minio.Client
}

type S3Options struct {
	Addr   string
	Key    string
	Secret string
	Secure bool
}

func NewS3(ctx context.Context, logger logging.Logger, opts S3Options) (*S3, error) {
	client, err := bootstrap.S3(ctx, bootstrap.S3Options(opts))
	if err != nil {
		return nil, fmt.Errorf("storage:S3.New: %w", err)
	}

	return &S3{
		logger: logger.With("system", "storage", "systemType", "s3"),
		client: client,
	}, nil
}

func (s *S3) Close(_ context.Context) error {
	return nil
}
