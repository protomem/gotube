package flashstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/protomem/gotube/internal/domain/model"
)

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

func (s *Storage) buildSessionKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}
