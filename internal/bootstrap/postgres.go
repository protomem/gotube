package bootstrap

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresOptions struct {
	Connect string
}

func Postgres(ctx context.Context, opts PostgresOptions) (*sql.DB, error) {
	const op = "bootstrap.Postgres"
	var err error

	db, err := sql.Open("pgx", opts.Connect)
	if err != nil {
		return nil, fmt.Errorf("%s: connect: %w", op, err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	return db, nil
}
