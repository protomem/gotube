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

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if err := run(logger); err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL  string
	httpPort int
	auth     struct {
		secret string
	}
	db struct {
		dsn         string
		automigrate bool
	}
	flash struct {
		dsn string
	}
	blob struct {
		addr   string
		key    string
		secret string
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

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:8080")
	cfg.httpPort = env.GetInt("HTTP_PORT", 8080)
	cfg.auth.secret = env.GetString("AUTH_SECRET", "secret")
	cfg.db.dsn = env.GetString("DB_DSN", "user:pass@localhost:5432/db")
	cfg.db.automigrate = env.GetBool("DB_AUTOMIGRATE", true)
	cfg.flash.dsn = env.GetString("FLASH_DSN", "localhost:6379/0")
	cfg.blob.addr = env.GetString("BLOB_ADDR", "localhost:9000")
	cfg.blob.key = env.GetString("BLOB_KEY", "user")
	cfg.blob.secret = env.GetString("BLOB_SECRET", "pass")

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

	fstore, err := flashstore.New(cfg.flash.dsn)
	if err != nil {
		return err
	}
	defer func() { _ = fstore.Close() }()

	bstore, err := blobstore.New(cfg.blob.addr, cfg.blob.key, cfg.blob.secret, false)
	if err != nil {
		return err
	}

	app := &application{
		config: cfg,
		db:     db,
		fstore: fstore,
		bstore: bstore,
		logger: logger,
	}

	return app.serveHTTP()
}
