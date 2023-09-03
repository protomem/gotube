package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ExpiredAt time.Time
	UserID    uuid.UUID
}

type SessionManager interface {
	GetSession(ctx context.Context, token string) (Session, error)
	SetSession(ctx context.Context, token string, sess Session) error
	DelSession(ctx context.Context, token string) error

	Close(ctx context.Context) error
}
