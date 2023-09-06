package storage

import (
	"context"
	"errors"
	"io"
)

var ErrObjectNotFound = errors.New("object not found")

type Object struct {
	Type string
	Size int64
	Data io.Reader
}

type BlobStorage interface {
	GetObject(ctx context.Context, parent string, name string) (Object, error)
	PutObject(ctx context.Context, parent string, name string, obj Object) error
	DelObject(ctx context.Context, parent string, name string) error

	Close(ctx context.Context) error
}
