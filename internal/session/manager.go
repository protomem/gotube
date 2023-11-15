package session

import "context"

type Manager interface {
	Get(ctx context.Context, token string) (Session, error)
	Set(ctx context.Context, sess Session) error
	Del(ctx context.Context, token string) error

	Close(ctx context.Context) error
}
