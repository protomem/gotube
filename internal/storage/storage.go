package storage

import "context"

type Storage interface {
	Close(ctx context.Context) error
}
