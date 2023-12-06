package session

import (
	"context"
	"fmt"

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

// TODO: implement it
func (r *Redis) Get(_ context.Context, token string) (Session, error) {
	return Session{}, nil
}

// TODO: implement it
func (r *Redis) Put(_ context.Context, sess Session) error {
	return nil
}

// TODO: implement it
func (r *Redis) Del(_ context.Context, token string) error {
	return nil
}

func (r *Redis) Close(_ context.Context) error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("session:Redis.Close: %w", err)
	}

	return nil
}
