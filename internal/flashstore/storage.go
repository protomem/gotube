package flashstore

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	_defaultTimeout = 3 * time.Second
	_defaultLeeway  = 15 * time.Second
)

type Storage struct {
	*redis.Client
}

func New(dsn string) (*Storage, error) {
	opts, err := redis.ParseURL("redis://" + dsn)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	return &Storage{client}, nil
}
