package flashstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token  string        `json:"token"`
	TTL    time.Duration `json:"ttl"`
	UserID uuid.UUID     `json:"userId"`
}

func (s *Storage) PutSession(ctx context.Context, session Session) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	sessBody, err := json.Marshal(session)
	if err != nil {
		return err
	}

	sessKey := fmt.Sprintf("session:%s", session.Token)

	if status := s.Set(ctx, sessKey, sessBody, session.TTL+_defaultLeeway); status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (s *Storage) DelSession(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	sessKey := fmt.Sprintf("session:%s", token)

	if status := s.Del(ctx, sessKey); status.Err() != nil {
		return status.Err()
	}

	return nil
}
