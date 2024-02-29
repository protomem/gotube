package database

import (
	"context"
	"io"
)

type Scanner interface {
	Scan(...any) error
}

type Rows interface {
	Scanner
	io.Closer

	Next() bool
}

type Row interface {
	Scanner
}

type DB interface {
	Exec(ctx context.Context, query string, args ...any) error
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row

	Close(ctx context.Context) error
}
