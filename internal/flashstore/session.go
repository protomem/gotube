package flashstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/protomem/gotube/internal/domain/model"
	"github.com/redis/go-redis/v9"
)

func (s *Storage) GetSession(ctx context.Context, token string) (model.Session, error) {
	const op = "flashstore.GetSession"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	sessionKey := s.buildSessionKey(token)

	res := s.Get(ctx, sessionKey)
	if res.Err() != nil {
		if res.Err() == redis.Nil {
			return model.Session{}, fmt.Errorf("%s: %w", op, model.ErrSessionNotFound)
		}

		return model.Session{}, fmt.Errorf("%s: %w", op, res.Err())
	}

	resData, err := res.Bytes()
	if err != nil {
		return model.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	var session model.Session
	if err := json.Unmarshal(resData, &session); err != nil {
		return model.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (s *Storage) PutSession(ctx context.Context, session model.Session) error {
	const op = "flashstore.PutSession"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	sessionKey := s.buildSessionKey(session.Token)

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res := s.Set(ctx, sessionKey, sessionJSON, time.Until(session.Expiry)); res.Err() != nil {
		return fmt.Errorf("%s: %w", op, res.Err())
	}

	return nil
}

func (s *Storage) DelSession(ctx context.Context, token string) error {
	const op = "flashstore.DelSession"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	sessionKey := s.buildSessionKey(token)

	if res := s.Del(ctx, sessionKey); res.Err() != nil {
		return fmt.Errorf("%s: %w", op, res.Err())
	}

	return nil
}

func (s *Storage) buildSessionKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}
