package flashstore

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	*redis.Client
}

func New(dsn string) (*Storage, error) {
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "flashstore.New", err)
	}

	client := redis.NewClient(opts)

	return &Storage{client}, nil
}
