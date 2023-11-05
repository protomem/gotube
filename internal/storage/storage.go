package storage

import (
	"context"
	"io"
)

type Object struct {
	Type string
	Size int64
	Src  io.Reader
}

type Storage interface {
	Get(ctx context.Context, parent, name string) (Object, error)
	Save(ctx context.Context, parent, name string, obj Object) error
	Delete(ctx context.Context, parent, name string) error

	Close(ctx context.Context) error
}
