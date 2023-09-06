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

func (*BlobStorage) GetObject(ctx context.Context, bucket string, object string) (storage.Object, error) {
	return storage.Object{}, nil
}

func (bs *BlobStorage) PutObject(ctx context.Context, bucket string, object string, src storage.Object) error {
	const op = "s3.BlobStorage.PutObject"
	var err error

	err = bs.initBucket(ctx, bucket)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = bs.client.PutObject(ctx, bucket, object, src.Data, src.Size, minio.PutObjectOptions{
		ContentType: src.Type,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (*BlobStorage) DelObject(ctx context.Context, bucket string, object string) error {
	return nil
}

func (*BlobStorage) Close(_ context.Context) error {
	return nil
}

func (bs *BlobStorage) initBucket(ctx context.Context, bucket string) error {
	const op = "init bucket"
	var err error

	exists, err := bs.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		err = bs.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
