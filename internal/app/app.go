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
	"github.com/protomem/gotube/internal/ctxstore"
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

	authConf, err := app.conf.Auth()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.repositories = sqliterepo.New(app.logger, app.db)
	app.services = service.New(authConf, app.repositories, bcrypt.New(bcrypt.DefaultCost))
	app.handlers = handler.New(app.logger, app.services)
	app.middlewares = middleware.New(app.logger, app.services)

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

	app.logger.Extractor(func(ctx context.Context) []any {
		args := []any{}
		if tid, ok := ctxstore.TraceID(ctx); ok {
			args = append(args, "traceId", tid)
		}
		return args
	})

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
	router := app.router
	middlewares := app.middlewares
	handlers := app.handlers

	router.Use(middlewares.TraceID())
	router.Use(middlewares.LogAccess(app.logger))
	router.Use(middlewares.Recovery(app.logger))

	router.Use(middlewares.Authenticate())

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowed)

	router.HandleFunc("/health", handlers.Health()).Methods(http.MethodGet)

	{
		router.HandleFunc("/users/{userNickname}", handlers.User.Get()).Methods(http.MethodGet)
		router.HandleFunc("/users", handlers.User.Create()).Methods(http.MethodPost)

		router.Handle(
			"/users/{userNickname}",
			middlewares.Protect()(handlers.User.Update()),
		).Methods(http.MethodPut, http.MethodPatch)
		router.Handle(
			"/users/{userNickname}",
			middlewares.Protect()(handlers.User.Delete()),
		).Methods(http.MethodDelete)

	}

	{
		router.Handle("/auth/login", handlers.Auth.Login()).Methods(http.MethodPost)
	}

	{
		router.HandleFunc("/subs/{userNickname}", handlers.Subscription.Count()).Methods(http.MethodGet)
		router.Handle(
			"/subs/{userNickname}",
			middlewares.Protect()(handlers.Subscription.Subscribe()),
		).Methods(http.MethodPost)
		router.Handle(
			"/subs/{userNickname}",
			middlewares.Protect()(handlers.Subscription.Unsubscribe()),
		).Methods(http.MethodDelete)
	}

	{
		router.HandleFunc("/videos", handlers.Video.List()).Methods(http.MethodGet)
		router.HandleFunc("/videos/{videoId}", handlers.Video.Get()).Methods(http.MethodGet)
		router.Handle(
			"/videos",
			middlewares.Protect()(handlers.Video.Creaate()),
		).Methods(http.MethodPost)
		router.Handle(
			"/videos/{videoId}",
			middlewares.Protect()(handlers.Video.Update()),
		).Methods(http.MethodPut, http.MethodPatch)
		router.Handle(
			"/videos/{videoId}",
			middlewares.Protect()(handlers.Video.Delete()),
		).Methods(http.MethodDelete)
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
