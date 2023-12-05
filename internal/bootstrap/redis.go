package bootstrap

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisOptions struct {
	Addr string

	Ping bool
}

func Redis(ctx context.Context, opts RedisOptions) (*redis.Client, error) {
	const op = "bootstrap.Redis"

	ropts := &redis.Options{
		Addr:     opts.Addr,
		Password: "",
		DB:       0,
	}

	client := redis.NewClient(ropts)

	if opts.Ping {
		if res := client.Ping(ctx); res.Err() != nil {
			return nil, fmt.Errorf("%s: ping: %w", op, res.Err())
		}
	}

	return client, nil
}
