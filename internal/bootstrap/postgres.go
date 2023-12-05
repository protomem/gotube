package bootstrap

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOptions struct {
	User     string
	Password string

	Host string
	Port int

	Database string

	Secure bool

	Ping bool
}

func Postgres(ctx context.Context, opts PostgresOptions) (*pgxpool.Pool, error) {
	const op = "bootstrap.Postgres"

	pool, err := pgxpool.New(
		ctx,
		buildPostgresConnect(opts.User, opts.Password, opts.Host, opts.Port, opts.Database, opts.Secure),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if opts.Ping {
		if err := pool.Ping(ctx); err != nil {
			return nil, fmt.Errorf("%s: ping: %w", op, err)
		}
	}

	return pool, nil
}

func buildPostgresConnect(user, password, host string, port int, database string, secure bool) string {
	connect := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, database)
	if secure {
		connect += "?sslmode=require"
	}
	return connect
}
