package blobstore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
)

var ErrObjectNotFound = errors.New("object not found")

type Storage interface {
	Get(ctx context.Context, parent, name string) (Object, error)
	Put(ctx context.Context, parent, name string, obj Object) error
	Del(ctx context.Context, parent, name string) error

	Close(ctx context.Context) error
}

type Object struct {
	Type string
	Size int64
	Body io.Reader
}

func (o Object) Clone() (Object, error) {
	buf := bytes.NewBuffer(make([]byte, 0, o.Size))
	if _, err := io.CopyN(buf, o.Body, o.Size); err != nil {
		return Object{}, fmt.Errorf("blobstore.Object: copy object: %w", err)
	}
	return Object{
		Type: o.Type,
		Size: o.Size,
		Body: buf,
	}, nil
}
