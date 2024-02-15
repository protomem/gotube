package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/blobstore"
	inmembstore "github.com/protomem/gotube/internal/blobstore/inmem"
	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/database"
	sqlitedb "github.com/protomem/gotube/internal/database/sqlite"
	"github.com/protomem/gotube/internal/handler"
	"github.com/protomem/gotube/internal/middleware"
	"github.com/protomem/gotube/internal/repository"
	sqliterepo "github.com/protomem/gotube/internal/repository/sqlite"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/closing"
	"github.com/protomem/gotube/pkg/hashing/bcrypt"
	"github.com/protomem/gotube/pkg/logging"
	stdlog "github.com/protomem/gotube/pkg/logging/std"
)

type App struct {
	conf   *config.Config
	logger logging.Logger

	db     database.DB
	bstore blobstore.Storage

	repositories *repository.Repositories
	services     *service.Services
	handlers     *handler.Handlers
	middlewares  *middleware.Middlewares

	router *mux.Router
	server *http.Server

	closer *closing.Closer
}

func New() *App {
	return &App{
		conf:   config.New(),
		router: mux.NewRouter(),
		closer: closing.New(),
	}
}

func (app *App) Run() error {
	const op = "app.Run"
	ctx := context.Background()

	if err := app.init(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.repositories = sqliterepo.New(app.logger, app.db)
	app.services = service.New(app.repositories, bcrypt.New(bcrypt.DefaultCost))
	app.handlers = handler.New(app.logger, app.services)
	app.middlewares = middleware.New()

	app.registerOnShutdown()
	app.setupRoutes()

	errs := make(chan error, 1)

	app.logger.Info("app initialized ...")
	defer app.logger.Info("app stopped.")

	go func() { app.serverStart(ctx, errs) }()
	go func() { app.gracefullShutdown(ctx, errs) }()

	if err := <-errs; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) init() error {
	const op = "init"
	ctx := context.Background()

	if err := app.initLogger(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := app.initDB(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := app.initBStore(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := app.initServer(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) initLogger() error {
	var err error

	conf, err := app.conf.Log()
	if err != nil {
		return err
	}

	app.logger, err = stdlog.New(conf.Level, os.Stdout)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) initDB(ctx context.Context) error {
	var err error

	conf, err := app.conf.SQLiteDB()
	if err != nil {
		return err
	}

	app.db, err = sqlitedb.Connect(ctx, app.logger, conf.DSN)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) initBStore() error {
	app.bstore, _ = inmembstore.New(app.logger)

	return nil
}

func (app *App) initServer() error {
	conf, err := app.conf.Server()
	if err != nil {
		return err
	}

	app.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler: app.router,
	}

	return nil
}

func (app *App) registerOnShutdown() {
	app.closer.Add(app.server.Shutdown)
	app.closer.Add(app.db.Close)
	app.closer.Add(app.bstore.Close)
	app.closer.Add(func(ctx context.Context) error {
		return app.logger.WithContext(ctx).Sync()
	})
}

func (app *App) setupRoutes() {
	app.router.Use(app.middlewares.TraceID())
	app.router.Use(app.middlewares.LogAccess(app.logger))
	app.router.Use(app.middlewares.Recovery(app.logger))

	app.router.NotFoundHandler = http.HandlerFunc(app.handlers.Common.NotFound)
	app.router.MethodNotAllowedHandler = http.HandlerFunc(app.handlers.MethodNotAllowed)

	app.router.HandleFunc("/health", app.handlers.Health()).Methods(http.MethodGet)

	{
		app.router.HandleFunc("/users/{userNickname}", app.handlers.User.Get()).Methods(http.MethodGet)
		app.router.HandleFunc("/users", app.handlers.User.Create()).Methods(http.MethodPost)
	}
}

func (app *App) serverStart(ctx context.Context, errs chan<- error) {
	app.logger.Info("app starting server ...", "addr", app.server.Addr)

	if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errs <- err
	}
}

func (app *App) gracefullShutdown(ctx context.Context, errs chan<- error) {
	<-wait()

	app.logger.Info("app shutting down ...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := app.closer.Close(ctx); err != nil {
		errs <- err
	}

	errs <- nil
}

func wait() <-chan os.Signal {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	return sigCh
}
