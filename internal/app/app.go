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

	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/pkg/closure"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/logging/stdlog"
)

type App struct {
	conf   config.Config
	logger logging.Logger

	router *http.ServeMux
	server *http.Server

	closer *closure.Closer
}

func New(conf config.Config) (*App, error) {
	const op = "app.New"
	var err error

	app := new(App)

	err = app.setup(conf)
	if err != nil {
		return nil, err
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

	app.conf = conf

	app.logger, err = stdlog.New(conf.Log.Level)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.router = http.NewServeMux()

	app.server = &http.Server{
		Addr:    conf.HTTP.Addr,
		Handler: app.router,
	}

	app.closer = closure.New()

	return nil
}

func (app *App) registerOnShutdown() {
	app.closer.Add(app.server.Shutdown)
	app.closer.Add(app.logger.Sync)
}

func (app *App) setupRoutes() {

	app.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GoTube v2.0")
	})

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
