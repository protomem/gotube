package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}
