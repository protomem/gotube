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

	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module"
	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/internal/storage/redis"
	"github.com/protomem/gotube/internal/storage/s3"
	"github.com/protomem/gotube/pkg/closer"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/logging/zap"
)

type App struct {
	conf   config.Config
	logger logging.Logger

	db *database.DB

	bstore  storage.BlobStorage
	sessmng storage.SessionManager

	modules *module.Modules

	app *echo.Echo

	closer *closer.Closer
}

func New(conf config.Config) (*App, error) {
	const op = "app.New"
	var err error
	ctx := context.Background()

	logger, err := zap.New(conf.Log.Level)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Debug("app configured ...", "conf", conf)

	db, err := database.New(ctx, logger, conf.Postgres.Connect)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	bstore, err := s3.NewBlobStorage(ctx, logger, s3.Options{
		Addr:      conf.S3.Addr,
		AccessKey: conf.S3.AccessKey,
		SecretKey: conf.S3.SecretKey,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sessmng, err := redis.NewSessionManager(ctx, logger, conf.Redis.Connect)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	modules := module.NewModules(logger, conf.Auth.Secret, db, bstore, sessmng)

	app := echo.New()
	app.HideBanner = true
	app.HidePort = true

	closer := closer.New()

	return &App{
		conf:    conf,
		logger:  logger,
		db:      db,
		bstore:  bstore,
		sessmng: sessmng,
		modules: modules,
		app:     app,
		closer:  closer,
	}, nil
}

func (app *App) Run() error {
	const op = "app.Run"
	var err error
	ctx := context.Background()

	app.registerOnShutdown()
	app.setupRoutes()

	errs := make(chan error)

	go app.startServer(ctx, errs)
	go app.gracefullShutdown(ctx, errs)

	app.logger.Info("app started ...", "addr", app.conf.HTTP.Addr)
	defer app.logger.Info("app stoped.")

	err = <-errs
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) registerOnShutdown() {
	app.closer.Add(app.app.Shutdown)
	app.closer.Add(app.sessmng.Close)
	app.closer.Add(app.bstore.Close)
	app.closer.Add(app.db.Close)
	app.closer.Add(app.logger.Sync)
}

func (app *App) setupRoutes() {
	app.router().Use(app.modules.Common.RequestID())
	app.router().Use(app.modules.Common.RequestLogger())
	app.router().Use(app.modules.Common.Recovery())
	app.router().Use(app.modules.Common.CORS())

	app.router().Use(app.modules.Auth.Authenticator())

	app.router().GET("/health", app.modules.Common.HandleHealthCheck())

	v1 := app.router().Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", app.modules.Auth.HandleRegister())
			auth.POST("/login", app.modules.Auth.HandleLogin())
			auth.DELETE("/logout", app.modules.Auth.HandleLogout(), app.modules.Auth.Authorizer())
			auth.GET("/refresh", app.modules.Auth.HandleRefreshTokens(), app.modules.Auth.Authorizer())
		}

		media := v1.Group("/media/files")
		{
			media.GET("/:folder/:file", app.modules.Media.HandleGetFile())
			media.POST("/:folder/:file", app.modules.Media.HandleSaveFile(), app.modules.Auth.Authorizer())
		}

		users := v1.Group("/users")
		{
			users.GET("/:nickname", app.modules.User.HandleGetUser())
			users.DELETE("/:nickname", app.modules.User.HandleDeleteUser(), app.modules.Auth.Authorizer())

			subscriptions := users.Group("/:nickname/subs")
			{
				subscriptions.Use(app.modules.Auth.Authorizer())

				subscriptions.GET("/", app.modules.Subscription.HandleGetAllSubscriptions())
				subscriptions.POST("/", app.modules.Subscription.HandleSubscribe())
				subscriptions.DELETE("/", app.modules.Subscription.HandleUnsubscribe())
			}
		}

		videos := v1.Group("/videos")
		{
			videos.GET("/:videoId", app.modules.Video.HandleGetVideo())
			videos.GET("/", app.modules.Video.HandleGetAllVideos())
			videos.POST("/", app.modules.Video.HandleCreateVideo(), app.modules.Auth.Authorizer())
		}
	}
}

func (app *App) startServer(_ context.Context, errs chan<- error) {
	err := app.app.Start(app.conf.HTTP.Addr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		errs <- fmt.Errorf("start server: %w", err)
	}
}

func (app *App) gracefullShutdown(ctx context.Context, errs chan<- error) {
	<-wait()

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err := app.closer.Close(ctx)
	if err != nil {
		errs <- fmt.Errorf("gracefull shutdown: %w", err)
	}

	errs <- nil
}

func (app *App) router() *echo.Echo {
	return app.app
}

func wait() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}
