package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/protomem/gotube/internal/bootstrap"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/redis/go-redis/v9"
)

var _ Manager = (*Redis)(nil)

type Redis struct {
	logger logging.Logger
	client *redis.Client
}

type RedisOptions struct {
	Addr string
	Ping bool
}

func NewRedis(ctx context.Context, logger logging.Logger, opts RedisOptions) (*Redis, error) {
	client, err := bootstrap.Redis(ctx, bootstrap.RedisOptions(opts))
	if err != nil {
		return nil, fmt.Errorf("session:Redis.New: %w", err)
	}

	return &Redis{
		logger: logger.With("system", "sessionManager", "systemType", "redis"),
		client: client,
	}, nil
}

func (r *Redis) Get(ctx context.Context, token string) (Session, error) {
	const op = "session:Redis.Get"

	res := r.client.Get(ctx, token)
	if res.Err() != nil {
		return Session{}, fmt.Errorf("%s: %w", op, res.Err())
	}

	val, err := res.Bytes()
	if err != nil {
		return Session{}, fmt.Errorf("%s: %w", op, err)
	}

	var sess Session
	if err := json.Unmarshal(val, &sess); err != nil {
		return Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return sess, nil
}

func (r *Redis) Put(ctx context.Context, sess Session) error {
	const op = "session:Redis.Put"

	sessJSON, err := json.Marshal(sess)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	sessTTL := time.Until(sess.Expiry)

	if res := r.client.Set(ctx, sess.Token, sessJSON, sessTTL); res.Err() != nil {
		return fmt.Errorf("%s: %w", op, res.Err())
	}

	return nil
}

func (r *Redis) Del(ctx context.Context, token string) error {
	const op = "session:Redis.Del"

	if res := r.client.Del(ctx, token); res.Err() != nil {
		return fmt.Errorf("%s: %w", op, res.Err())
	}

	return nil
}

func (r *Redis) Close(_ context.Context) error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("session:Redis.Close: %w", err)
	}

	return nil
}
