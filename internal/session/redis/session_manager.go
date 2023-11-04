package redis

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/bootstrap"
	"github.com/protomem/gotube/internal/session"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/redis/go-redis/v9"
)

var _ session.Manager = (*SessionManager)(nil)

type SessionManager struct {
	logger logging.Logger
	rdb    *redis.Client
}

func NewSessionManager(ctx context.Context, logger logging.Logger, addr string) (*SessionManager, error) {
	const op = "redis.NewSessionManager"

	rdb, err := bootstrap.Redis(ctx, bootstrap.RedisOptions{
		Addr: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &SessionManager{
		logger: logger.With("component", "sessionManager", "sessionManagerType", "redis"),
		rdb:    rdb,
	}, nil
}

func (mng *SessionManager) Get(ctx context.Context, token string) (session.Session, error) {
	const op = "redis.SessionManager.Get"

	res := mng.rdb.HGetAll(ctx, token)
	if res.Err() != nil {
		return session.Session{}, fmt.Errorf("%s: %w", op, res.Err())
	}

	var sess session.Session
	err := res.Scan(&sess)
	if err != nil {
		return session.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return sess, nil
}

func (mng *SessionManager) Set(ctx context.Context, session session.Session) error {
	const op = "redis.SessionManager.Set"

	res := mng.rdb.HSet(ctx, session.Token, session)
	if res.Err() != nil {
		return fmt.Errorf("%s: %w", op, res.Err())
	}

	return nil
}

func (mng *SessionManager) Del(ctx context.Context, token string) error {
	const op = "redis.SessionManager.Del"

	res := mng.rdb.Del(ctx, token)
	if res.Err() != nil {
		return fmt.Errorf("%s: %w", op, res.Err())
	}

	return nil
}

func (mng *SessionManager) Close(_ context.Context) error {
	const op = "redis.SessionManager.Close"

	err := mng.rdb.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
