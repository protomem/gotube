package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protomem/gotube/pkg/logging"
)

type DB struct {
	logger logging.Logger
	pool   *pgxpool.Pool
}

func New(ctx context.Context, logger logging.Logger, connect string) (*DB, error) {
	const op = "db.New"
	var err error

	pool, err := pgxpool.New(ctx, connect)
	if err != nil {
		return nil, fmt.Errorf("%s: connect: %w", op, err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	return &DB{
		logger: logger.With("system", "db", "dbType", "postgres"),
		pool:   pool,
	}, nil
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return db.pool.QueryRow(ctx, query, args...)
}

func (db *DB) Close(_ context.Context) error {
	db.pool.Close()
	return nil
}
