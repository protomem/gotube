package redis

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/redis/go-redis/v9"
)

var _ storage.SessionManager = (*SessionManager)(nil)

type SessionManager struct {
	logger logging.Logger
	client *redis.Client
}

func NewSessionManager(ctx context.Context, logger logging.Logger, connect string) (*SessionManager, error) {
	const op = "redis.SessionManager.New"

	opts, err := redis.ParseURL(connect)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := redis.NewClient(opts)

	return &SessionManager{
		logger: logger.With("system", "sessionManager", "sessionManagerType", "redis"),
		client: client,
	}, nil
}

func (sm *SessionManager) GetSession(ctx context.Context, token string) (storage.Session, error) {
	const op = "sessionManager.GetSession"
	var err error

	res := sm.client.HGetAll(ctx, token)
	if res.Err() != nil {
		return storage.Session{}, fmt.Errorf("%s: %w", op, res.Err())
	}

	var sess storage.Session
	err = res.Scan(&sess)
	if err != nil {
		return storage.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return sess, nil
}

func (sm *SessionManager) SetSession(ctx context.Context, token string, sess storage.Session) error {
	const op = "sessionManager.SetSession"

	res := sm.client.HSet(ctx, token, sess)
	if res.Err() != nil {
		return fmt.Errorf("%s: %w", op, res.Err())
	}

	return nil
}

func (sm *SessionManager) DelSession(ctx context.Context, token string) error {
	const op = "sessionManager.DelSession"

	res := sm.client.Del(ctx, token)
	if res.Err() != nil {
		return fmt.Errorf("%s: %w", op, res.Err())
	}

	return nil
}

func (*SessionManager) Close(_ context.Context) error {
	return nil
}
