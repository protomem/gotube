package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/bootstrap"
	"github.com/protomem/gotube/internal/config"
	httphandl "github.com/protomem/gotube/internal/handler/http"
	"github.com/protomem/gotube/internal/repository"
	postgresrepo "github.com/protomem/gotube/internal/repository/postgres"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/internal/session"
	"github.com/protomem/gotube/internal/session/redis"
	"github.com/protomem/gotube/pkg/closure"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/logging/stdlog"
)

type repositories struct {
	repository.User
}

func newRepositories(logger logging.Logger, pgdb *sql.DB) *repositories {
	return &repositories{
		User: postgresrepo.NewUserRepository(logger, pgdb),
	}
}

type services struct {
	service.User
	service.Auth
}

func newServices(authSigngingKey string, repos *repositories, sessmng session.Manager) *services {
	servs := &services{}

	servs.User = service.NewUser(repos.User)
	servs.Auth = service.NewAuth(authSigngingKey, sessmng, servs.User)

	return servs
}

type handlers struct {
	*httphandl.UserHandler
	*httphandl.AuthHandler
}

func newHandlers(logger logging.Logger, services *services) *handlers {
	return &handlers{
		UserHandler: httphandl.NewUserHandler(logger, services.User),
		AuthHandler: httphandl.NewAuthHandler(logger, services.Auth),
	}
}

type App struct {
	conf   config.Config
	logger logging.Logger

	pgdb *sql.DB

	sessmng session.Manager

	repos  *repositories
	servs  *services
	handls *handlers

	router *mux.Router
	server *http.Server

	closer *closure.Closer
}

func New(conf config.Config) (*App, error) {
	const op = "app.New"
	var err error

	app := new(App)

	err = app.setup(conf)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (app *App) Run() error {
	const op = "app.Run"
	var err error

	app.registerOnShutdown()
	app.setupRoutes()

	app.logger.Debug("app configured", "config", app.conf)

	errs := make(chan error)

	go func() { app.startServer(errs) }()
	go func() { app.gracefulShutdown(errs) }()

	app.logger.Info("app started ...")
	defer app.logger.Info("app finished")

	err = <-errs
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) setup(conf config.Config) error {
	const op = "setup"
	var err error
	ctx := context.Background()

	app.conf = conf

	app.logger, err = stdlog.New(conf.Log.Level)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.pgdb, err = bootstrap.Postgres(ctx, bootstrap.PostgresOptions{
		Connect: conf.Postgres.Connect,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.sessmng, err = redis.NewSessionManager(ctx, app.logger, conf.Redis.Addr)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.repos = newRepositories(app.logger, app.pgdb)
	app.servs = newServices(conf.Auth.Secret, app.repos, app.sessmng)
	app.handls = newHandlers(app.logger, app.servs)

	app.router = mux.NewRouter()

	app.server = &http.Server{
		Addr:    conf.HTTP.Addr,
		Handler: app.router,
	}

	app.closer = closure.New()

	return nil
}

func (app *App) registerOnShutdown() {
	app.closer.Add(app.server.Shutdown)
	app.closer.Add(app.sessmng.Close)
	app.closer.Add(func(_ context.Context) error {
		return app.pgdb.Close()
	})
	app.closer.Add(app.logger.Sync)
}

func (app *App) setupRoutes() {

	app.router.HandleFunc("/api/v1/auth/register", app.handls.AuthHandler.Register()).Methods(http.MethodPost)
	app.router.HandleFunc("/api/v1/auth/login", app.handls.AuthHandler.Login()).Methods(http.MethodPost)
	app.router.HandleFunc("/api/v1/auth/refresh", app.handls.AuthHandler.Refresh()).Methods(http.MethodPost)
	app.router.HandleFunc("/api/v1/auth/logout", app.handls.AuthHandler.Logout()).Methods(http.MethodPost)

	app.router.HandleFunc("/api/v1/users/{nickname}", app.handls.UserHandler.Get()).Methods(http.MethodGet)
	app.router.HandleFunc("/api/v1/users/{nickname}", app.handls.UserHandler.Update()).Methods(http.MethodPatch)
	app.router.HandleFunc("/api/v1/users/{nickname}", app.handls.UserHandler.Delete()).Methods(http.MethodDelete)

}

func (app *App) startServer(errs chan<- error) {
	const op = "startServer"

	err := app.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		errs <- fmt.Errorf("%s: %w", op, err)
	}
}

func (app *App) gracefulShutdown(errs chan<- error) {
	const op = "gracefulShutdown"

	<-wait()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := app.closer.Close(ctx)
	if err != nil {
		errs <- fmt.Errorf("%s: %w", op, err)
	}

	errs <- nil
}

func wait() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}
