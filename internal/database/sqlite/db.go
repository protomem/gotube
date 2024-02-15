package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/pkg/logging"
)

var _ database.DB = (*DB)(nil)

type DB struct {
	*sql.DB

	logger logging.Logger
}

func Connect(ctx context.Context, logger logging.Logger, dsn string) (*DB, error) {
	const op = "database.Connect"

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &DB{
		DB:     db,
		logger: logger.With("component", "sqlite/database"),
	}, nil
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) error {
	const op = "database.Exec"
	db.logger.WithContext(ctx).Debug("query exec", "query", query, "args", args, "operation", op)

	if _, err := db.DB.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (db *DB) Query(ctx context.Context, query string, args ...any) (database.Rows, error) {
	const op = "database.Query"
	db.logger.WithContext(ctx).Debug("query exec", "query", query, "args", args, "operation", op)

	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return rows, nil
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...any) database.Row {
	const op = "database.QueryRow"
	db.logger.WithContext(ctx).Debug("query exec", "query", query, "args", args, "operation", op)

	return db.DB.QueryRowContext(ctx, query, args...)
}

func (db *DB) Close(_ context.Context) error {
	if err := db.DB.Close(); err != nil {
		return fmt.Errorf("database.Close: %w", err)
	}
	return nil
}
