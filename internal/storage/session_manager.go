package storage

import (
	"context"
	"time"
)

type Session struct {
	ExpiredAt time.Time `redis:"expiredAt"`
	UserID    string    `redis:"userId"`
}

type SessionManager interface {
	GetSession(ctx context.Context, token string) (Session, error)
	SetSession(ctx context.Context, token string, sess Session) error
	DelSession(ctx context.Context, token string) error

	Close(ctx context.Context) error
}
