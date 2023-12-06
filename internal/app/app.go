package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protomem/gotube/internal/bootstrap"
	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/handler"
	"github.com/protomem/gotube/internal/middleware"
	"github.com/protomem/gotube/pkg/closing"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/logging/zap"
	"github.com/protomem/gotube/pkg/requestid"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	conf   config.Config
	logger logging.Logger

	pdb *pgxpool.Pool
	mdb *mongo.Client

	handls *handler.Handlers
	mdws   *middleware.Middlewares

	router *chi.Mux
	server *http.Server

	closer *closing.Closer
}

func New(conf config.Config) *App {
	return &App{conf: conf}
}

func (app *App) Run(ctx context.Context) error {
	const op = "application:App.Run"

	if err := app.setup(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.handls = handler.New(app.logger)
	app.mdws = middleware.New(app.logger)

	app.registerOnShutdown()
	app.setupRoutes()

	errs := make(chan error, 1)

	go func() { app.startServer(errs) }()
	go func() { app.gracefulShutdown(errs) }()

	app.logger.Info("app started ...", "addr", app.server.Addr)
	defer app.logger.Info("app stopped.")

	if err := <-errs; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) setup(ctx context.Context) error {
	const op = "setup"
	var err error

	if app.logger, err = zap.NewLogger(app.conf.Log.Level); err != nil {
		return fmt.Errorf("%s: logger: %w", op, err)
	}

	app.logger.Debug("app configure ...", "config", app.conf)

	if app.pdb, err = bootstrap.Postgres(ctx, bootstrap.PostgresOptions{
		Host:     app.conf.Postgres.Host,
		Port:     app.conf.Postgres.Port,
		User:     app.conf.Postgres.User,
		Password: app.conf.Postgres.Password,
		Database: app.conf.Postgres.Database,
		Secure:   app.conf.Postgres.Secure,
		Ping:     app.conf.Mode == config.Dev,
	}); err != nil {
		return fmt.Errorf("%s: postgres: %w", op, err)
	}

	if app.mdb, err = bootstrap.Mongo(ctx, bootstrap.MongoOptions{
		URI:  app.conf.Mongo.URI,
		Ping: app.conf.Mode == config.Dev,
	}); err != nil {
		return fmt.Errorf("%s: mongo: %w", op, err)
	}

	app.router = chi.NewRouter()
	app.server = &http.Server{
		Addr:    httpAddr(app.conf.HTTP.Host, app.conf.HTTP.Port),
		Handler: app.router,
	}

	app.closer = closing.New()

	return nil
}

func (app *App) registerOnShutdown() {
	app.closer.Add(app.server.Shutdown)
	app.closer.Add(app.mdb.Disconnect)
	app.closer.Add(func(_ context.Context) error {
		app.pdb.Close()
		return nil
	})
	app.closer.Add(app.logger.Sync)
}

func (app *App) setupRoutes() {
	app.router.Use(app.mdws.CORS())
	app.router.Use(requestid.Middleware())
	app.router.Use(app.mdws.RealIP())
	app.router.Use(app.mdws.CleanPath())
	app.router.Use(app.mdws.StripSlashes())
	app.router.Use(app.mdws.RequestLogging())
	app.router.Use(app.mdws.Recoverer())

	app.router.Get("/ping", app.handls.Ping())
}

func (app *App) startServer(errs chan<- error) {
	if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errs <- fmt.Errorf("start server: %w", err)
		return
	}
}

func (app *App) gracefulShutdown(errs chan<- error) {
	<-wait()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.closer.Close(ctx); err != nil {
		errs <- fmt.Errorf("graceful shutdown: %w", err)
		return
	}

	errs <- nil
}

func wait() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}

func httpAddr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
