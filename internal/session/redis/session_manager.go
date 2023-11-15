package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/bootstrap"
	"github.com/protomem/gotube/internal/session"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/redis/go-redis/v9"
)

var _ session.Manager = (*SessionManager)(nil)

type customSession struct {
	UserId    string `redis:"userId"`
	ExpiresAt string `redis:"expiresAt"`
}

func (sess customSession) ToSession(token string) (session.Session, error) {
	var err error

	userId, err := uuid.Parse(sess.UserId)
	if err != nil {
		return session.Session{}, err
	}

	expiresAt, err := time.Parse(time.RFC3339, sess.ExpiresAt)
	if err != nil {
		return session.Session{}, err
	}

	return session.Session{
		Token:     token,
		UserID:    userId,
		ExpiresAt: expiresAt,
	}, nil
}

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

	var csess customSession
	err := res.Scan(&csess)
	if err != nil {
		return session.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	sess, err := csess.ToSession(token)
	if err != nil {
		return session.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return sess, nil
}

func (mng *SessionManager) Set(ctx context.Context, sess session.Session) error {
	const op = "redis.SessionManager.Set"

	csess := customSession{
		UserId:    sess.UserID.String(),
		ExpiresAt: sess.ExpiresAt.Format(time.RFC3339),
	}

	res := mng.rdb.HSet(ctx, sess.Token, csess)
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
