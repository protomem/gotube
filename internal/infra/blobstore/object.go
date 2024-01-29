package blobstore

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

var ErrObjectNotFound = errors.New("object not found")

type Object struct {
	Type string
	Size int64
	Body io.ReadCloser
}

func (s *Storage) GetObject(ctx context.Context, bucketName, objectName string) (Object, error) {
	const op = "blobstore.GetObject"

	opts := minio.GetObjectOptions{}

	object, err := s.Client.GetObject(ctx, bucketName, objectName, opts)
	if err != nil {
		return Object{}, fmt.Errorf("%s: %w", op, err)
	}

	objectInfo, err := object.Stat()
	if err != nil {
		return Object{}, fmt.Errorf("%s: %w", op, err)
	}

	return Object{
		Type: objectInfo.ContentType,
		Size: objectInfo.Size,
		Body: object,
	}, nil
}

func (s *Storage) PutObject(ctx context.Context, bucketName, objectName string, object Object) error {
	const op = "blobstore.PutObject"

	if err := s.initBucket(ctx, bucketName); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	opts := minio.PutObjectOptions{
		ContentType: object.Type,
	}

	if _, err := s.Client.PutObject(ctx, bucketName, objectName, object.Body, object.Size, opts); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	const op = "blobstore.DeleteObject"

	opts := minio.RemoveObjectOptions{}

	if err := s.Client.RemoveObject(ctx, bucketName, objectName, opts); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) initBucket(ctx context.Context, bucketName string) error {
	const op = "initBucket"

	exists, err := s.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if exists {
		return nil
	}

	opts := minio.MakeBucketOptions{}

	if err := s.MakeBucket(ctx, bucketName, opts); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
