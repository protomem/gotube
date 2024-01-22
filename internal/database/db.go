package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/protomem/gotube/assets"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const _defaultTimeout = 3 * time.Second

type DB struct {
	*sqlx.DB
}

func New(dsn string, automigrate bool) (*DB, error) {
	const op = "database.New"

	ctx, cancel := context.WithTimeout(context.Background(), _defaultTimeout)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "pgx", "postgres://"+dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	if automigrate {
		const subOp = op + ": automigrate"

		iofsDriver, err := iofs.New(assets.Assets, "migrations")
		if err != nil {
			return nil, fmt.Errorf("%s: %w", subOp, err)
		}

		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "postgres://"+dsn)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", subOp, err)
		}

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, fmt.Errorf("%s: %w", subOp, err)
		}
	}

	return &DB{db}, nil
}
