package database

import "context"

type Rows interface{}

type Row interface{}

type DB interface {
	Exec(ctx context.Context, query string, args ...any) error
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row

	Close(ctx context.Context) error
}
