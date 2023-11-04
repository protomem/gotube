package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token     string    `redis:"-"`
	UserID    uuid.UUID `redis:"userId"`
	ExpiresAt time.Time `redis:"expiresAt"`
}
