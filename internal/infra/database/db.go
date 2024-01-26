package database

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/protomem/gotube/internal/config"
)

const _defaultTimeout = 5 * time.Second

type DB struct {
	*sqlx.DB

	conf config.Database
}

func New(conf config.Config) (*DB, error) {
	return &DB{conf: conf.Database}, nil
}

func (db *DB) Connect(ctx context.Context) error {
	var err error

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	db.DB, err = sqlx.ConnectContext(ctx, "pgx", "postgres://"+db.conf.DSN)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	return nil
}

func (db *DB) Disconnect(ctx context.Context) error {
	return db.DB.Close()
}

type SelectOptions struct {
	Limit  uint64
	Offset uint64
}
