package application

import (
	"os"

	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/infra/database"
	"github.com/protomem/gotube/internal/infra/flashstore"
	"github.com/protomem/gotube/internal/infra/routes"
	"github.com/protomem/gotube/internal/infra/server"
	"github.com/protomem/gotube/pkg/logging"
	"go.uber.org/fx"
)

func ProvideLogger(conf config.Config) (*logging.Logger, error) {
	return logging.New(os.Stdout, conf.Log.Level)
}

func Create() fx.Option {
	return fx.Options(
		fx.Provide(
			config.New,
			ProvideLogger,
			database.New,
			flashstore.New,
			routes.Setup,
			server.New,
		),
		fx.Invoke(
			registerRunners,
		),
	)
}

func registerRunners(lc fx.Lifecycle, db *database.DB, fstore *flashstore.Storage, srv *server.Server) {
	lc.Append(fx.Hook{
		OnStart: db.Connect,
		OnStop:  db.Disconnect,
	})

	lc.Append(fx.Hook{
		OnStart: fstore.Connect,
		OnStop:  fstore.Disconnect,
	})

	lc.Append(fx.Hook{
		OnStart: srv.Start,
		OnStop:  srv.Stop,
	})
}
