package main

import (
	"errors"
	"flag"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/protomem/gotube/assets"
	"github.com/protomem/gotube/internal/env"
)

func main() {
	migrationAction := flag.String("action", "up", "migration action")
	migrationStep := flag.Uint("step", 0, "migration step")

	db := flag.String("db", "", "database dsn")
	if *db == "" {
		*db = env.GetString("DB_DSN", "")
		if *db == "" {
			panic("database dsn is required")
		}
	}

	flag.Parse()

	iofsDriver, err := iofs.New(assets.Assetss, "migrations")
	if err != nil {
		panic(err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "postgres://"+*db)
	if err != nil {
		panic(err)
	}

	switch *migrationAction {
	case "up":
		if *migrationStep != 0 {
			if err := migrator.Steps(int(*migrationStep)); err != nil {
				panic(err)
			}
		} else {
			if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				panic(err)
			}
		}
	case "down":
		if *migrationStep != 0 {
			if err := migrator.Steps(int(*migrationStep)); err != nil {
				panic(err)
			}
		} else {
			if err := migrator.Down(); err != nil {
				panic(err)
			}
		}
	case "drop":
		if err := migrator.Drop(); err != nil {
			panic(err)
		}
	default:
		panic("unknown action, supported actions: up [step], down [step], drop")
	}
}
