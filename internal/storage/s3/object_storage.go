package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/protomem/gotube/internal/bootstrap"
	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/logging"
)

var _ storage.Storage = (*ObjectStorage)(nil)

type ObjectStorage struct {
	logger logging.Logger
	s3db   *minio.Client
}

func NewObjectStorage(ctx context.Context, logger logging.Logger, addr, access, secret string) (*ObjectStorage, error) {
	const op = "s3.NewObjectStorage"

	s3db, err := bootstrap.S3(ctx, bootstrap.S3Options{
		Addr:   addr,
		Access: access,
		Secret: secret,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &ObjectStorage{
		logger: logger.With("component", "objectStorage", "objectStorageType", "s3"),
		s3db:   s3db,
	}, nil
}

func (s *ObjectStorage) Get(ctx context.Context, parent, name string) (storage.Object, error) {
	panic("unimplemented")
}

func (s *ObjectStorage) Save(ctx context.Context, parent, name string, obj storage.Object) error {
	panic("unimplemented")
}

func (s *ObjectStorage) Delete(ctx context.Context, parent, nama string) error {
	panic("unimplemented")
}

func (s *ObjectStorage) Close(_ context.Context) error {
	return nil
}
