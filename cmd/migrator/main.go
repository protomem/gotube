package main

import (
	"errors"
	"flag"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/protomem/gotube/assets"
)

const (
	_migrateUp   = "up"
	_migrateDown = "down"
	_migrateDrop = "drop"
	_migrateStep = "step"
)

var (
	_databaseURL     = flag.String("database", "", "database url")
	_migrationAction = flag.String("action", "up", "migrations action")
	_migrationStep   = flag.String("step", "", "migrations step")
)

func init() {
	flag.Parse()
}

func main() {
	var err error

	if *_databaseURL == "" {
		panic("database url not set")
	}

	source, err := iofs.New(assets.Assets, "migrations")
	if err != nil {
		panic("failed to create migrations source: " + err.Error())
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, *_databaseURL)
	if err != nil {
		panic("failed to create migrations instance: " + err.Error())
	}
	defer func() { _, _ = m.Close() }()

	// TODO: add step flag
	switch *_migrationAction {
	case _migrateUp:
		err = m.Up()
	case _migrateDown:
		err = m.Down()
	case _migrateDrop:
		err = m.Drop()
	default:
		panic("migrations action must be 'up' or 'down' or 'drop'")
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic("failed to run " + *_migrationAction + " migrations: " + err.Error())
	}
}
