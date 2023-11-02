package bootstrap

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisOptions struct {
	Addr string
}

func Redis(ctx context.Context, opts RedisOptions) (*redis.Client, error) {
	const op = "bootstrap.Redis"

	client := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	res := client.Ping(ctx)
	if res.Err() != nil {
		return nil, fmt.Errorf("%s: ping: %w", op, res.Err())
	}

	return client, nil
}
