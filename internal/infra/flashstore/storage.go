package flashstore

import (
	"context"

	"github.com/protomem/gotube/internal/config"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	*redis.Client

	conf config.Flash
}

func New(conf config.Config) *Storage {
	return &Storage{conf: conf.Flash}
}

func (s *Storage) Connect(_ context.Context) error {
	var err error

	opts, err := redis.ParseURL("redis://" + s.conf.DSN)
	if err != nil {
		return err
	}

	s.Client = redis.NewClient(opts)

	return nil
}

func (s *Storage) Disconnect(ctx context.Context) error {
	return s.Client.Close()
}
