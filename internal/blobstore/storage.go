package blobstore

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var ErrObjectNotFound = errors.New("object not found")

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

type Object struct {
	Type string
	Size int64
	Body io.ReadCloser
}

func (s *Storage) GetFile(ctx context.Context, bucket, name string) (Object, error) {
	opts := minio.GetObjectOptions{}

	obj, err := s.GetObject(ctx, bucket, name, opts)
	if err != nil {
		return Object{}, err
	}

	objStat, err := obj.Stat()
	if err != nil {
		res := minio.ToErrorResponse(err)
		if res.StatusCode == http.StatusNotFound {
			return Object{}, ErrObjectNotFound
		}

		return Object{}, err
	}

	return Object{
		Type: objStat.ContentType,
		Size: objStat.Size,
		Body: obj,
	}, nil
}

func (s *Storage) SaveFile(ctx context.Context, bucket, name string, obj Object) error {
	if err := s.initBucket(ctx, bucket); err != nil {
		return err
	}

	opts := minio.PutObjectOptions{
		ContentType:           obj.Type,
		ConcurrentStreamParts: true,
		NumThreads:            4,
	}

	if _, err := s.PutObject(ctx, bucket, name, obj.Body, obj.Size, opts); err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveFile(ctx context.Context, bucket, name string) error {
	opts := minio.RemoveObjectOptions{}

	if err := s.RemoveObject(ctx, bucket, name, opts); err != nil {
		return err
	}

	return nil
}

func (s *Storage) initBucket(ctx context.Context, bucket string) error {
	if exists, err := s.BucketExists(ctx, bucket); err != nil || exists {
		if exists {
			return nil
		}

		return err
	}

	opts := minio.MakeBucketOptions{}

	if err := s.MakeBucket(ctx, bucket, opts); err != nil {
		return err
	}

	return nil
}
