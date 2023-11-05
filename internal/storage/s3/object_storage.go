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

func (s *ObjectStorage) Get(ctx context.Context, bucketName, objectName string) (storage.Object, error) {
	const op = "s3.ObjectStorage.Get"
	var err error

	// TODO: Add handle for object not found

	obj, err := s.s3db.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return storage.Object{}, fmt.Errorf("%s: %w", op, err)
	}

	objInfo, err := obj.Stat()
	if err != nil {
		return storage.Object{}, fmt.Errorf("%s: %w", op, err)
	}

	return storage.Object{
		Type: objInfo.ContentType,
		Size: objInfo.Size,
		Src:  obj,
	}, nil
}

func (s *ObjectStorage) Save(ctx context.Context, bucketName, objectName string, obj storage.Object) error {
	const op = "s3.ObjectStorage.Save"
	var err error

	err = s.initBucket(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.s3db.PutObject(ctx, bucketName, objectName, obj.Src, obj.Size, minio.PutObjectOptions{
		ContentType: obj.Type,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *ObjectStorage) Delete(ctx context.Context, bucketName, objectName string) error {
	const op = "s3.ObjectStorage.Delete"

	err := s.s3db.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *ObjectStorage) Close(_ context.Context) error {
	return nil
}

func (s *ObjectStorage) initBucket(ctx context.Context) error {
	return nil
}
