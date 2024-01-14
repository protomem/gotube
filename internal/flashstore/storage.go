package flashstore

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const _defaultTimeout = 2 * time.Second

type Storage struct {
	*redis.Client
}

func New(dsn string) (*Storage, error) {
	opts, err := redis.ParseURL("redis://" + dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "flashstore.New", err)
	}

	client := redis.NewClient(opts)

	return &Storage{client}, nil
}
