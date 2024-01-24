package flashstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/protomem/gotube/internal/config"
	"github.com/redis/go-redis/v9"
)

var ErrKeyNotFound = errors.New("key not found")

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

func (s *Storage) ScanJSON(ctx context.Context, key string, value any) error {
	const op = "flashstore.GetJSON"

	valueJSON, err := s.Client.Get(ctx, key).Bytes()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := json.Unmarshal(valueJSON, value); err != nil {
		if err == redis.Nil {
			return fmt.Errorf("%s: %w", op, ErrKeyNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) PutJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	const op = "flashstore.PutJSON"

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.Client.Set(ctx, key, valueJSON, ttl).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	const op = "flashstore.Delete"

	if err := s.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
