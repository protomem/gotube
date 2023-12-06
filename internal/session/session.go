package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Duration
}

type Manager interface {
	Get(ctx context.Context, token string) (Session, error)
	Put(ctx context.Context, session Session) error
	Del(ctx context.Context, token string) error

	Close(ctx context.Context) error
}
