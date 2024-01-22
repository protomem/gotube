package application

import (
	"net/http"

	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/infra/routes"
	"github.com/protomem/gotube/internal/infra/server"
	"go.uber.org/fx"
)

func Create() fx.Option {
	return fx.Options(
		fx.Provide(
			config.New,
			fx.Annotate(routes.Setup, fx.As(new(http.Handler))),
			server.New,
		),
		fx.Invoke(
			registerRunners,
		),
	)
}

func registerRunners(lc fx.Lifecycle, srv *server.Server) {
	lc.Append(fx.Hook{
		OnStart: srv.Start,
		OnStop:  srv.Stop,
	})
}
