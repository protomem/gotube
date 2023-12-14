package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/protomem/gotube/internal/blobstore"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/env"
	"github.com/protomem/gotube/internal/flashstore"
	"github.com/protomem/gotube/internal/version"
)

// TODO: Add request id
// TODO: Add docker, docker-compose
// TODO: Add JWT
// TODO: Add access and refresh tokens
// TODO: Add Casbin

// TODO: Add subscriptions
// TODO: Add videos
// TODO: Add raitings
// TODO: Add comments
// TODO: Add media

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL  string
	httpPort int
	cookie   struct {
		secretKey string
	}
	db struct {
		dsn         string
		automigrate bool
	}
	flash struct {
		dsn string
	}
	blob struct {
		addr      string
		accessKey string
		secretKey string
	}
}

type application struct {
	config config
	db     *database.DB
	fstore *flashstore.Storage
	bstore *blobstore.Storage
	logger *slog.Logger
	wg     sync.WaitGroup
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:4444")
	cfg.httpPort = env.GetInt("HTTP_PORT", 4444)
	cfg.cookie.secretKey = env.GetString("COOKIE_SECRET_KEY", "34gqafpamolgom4njk6wjcxh2qilxdwd")
	cfg.db.dsn = env.GetString("DB_DSN", "user:pass@localhost:5432/db")
	cfg.db.automigrate = env.GetBool("DB_AUTOMIGRATE", true)
	cfg.flash.dsn = env.GetString("FLASH_DSN", "localhost:6379/0")
	cfg.blob.addr = env.GetString("BLOB_ADDR", "localhost:9000")
	cfg.blob.accessKey = env.GetString("BLOB_ACCESS_KEY", "user")
	cfg.blob.secretKey = env.GetString("BLOB_SECRET_KEY", "pass")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	db, err := database.New(cfg.db.dsn, cfg.db.automigrate)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	bstore, err := blobstore.New(cfg.blob.addr, cfg.blob.accessKey, cfg.blob.secretKey, false)
	if err != nil {
		return err
	}

	fstore, err := flashstore.New(cfg.flash.dsn)
	if err != nil {
		return err
	}
	defer func() { _ = fstore.Close() }()

	app := &application{
		config: cfg,
		db:     db,
		bstore: bstore,
		fstore: fstore,
		logger: logger,
	}

	return app.serveHTTP()
}
