package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/logging"
)

var _ storage.BlobStorage = (*BlobStorage)(nil)

type BlobStorage struct {
	logger logging.Logger
	client *minio.Client
}

type Options struct {
	Addr      string
	AccessKey string
	SecretKey string
}

func NewBlobStorage(ctx context.Context, logger logging.Logger, opts Options) (*BlobStorage, error) {
	const op = "s3.BlobStorage.New"

	client, err := minio.New(opts.Addr, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.AccessKey, opts.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &BlobStorage{
		logger: logger.With("system", "blobStorage", "blobStorageType", "s3"),
		client: client,
	}, nil
}

func (*BlobStorage) GetObject(ctx context.Context, parent string, name string) (storage.Object, error) {
	return storage.Object{}, nil
}

func (*BlobStorage) PutObject(ctx context.Context, parent string, name string, obj storage.Object) error {
	return nil
}

func (*BlobStorage) DelObject(ctx context.Context, parent string, name string) error {
	return nil
}

func (*BlobStorage) Close(_ context.Context) error {
	return nil
}
