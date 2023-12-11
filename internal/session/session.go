package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"userId"`
	Expiry time.Time `json:"expiry"`
}

type Manager interface {
	Get(ctx context.Context, token string) (Session, error)
	Put(ctx context.Context, session Session) error
	Del(ctx context.Context, token string) error

	Close(ctx context.Context) error
}
