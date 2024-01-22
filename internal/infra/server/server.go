package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/protomem/gotube/internal/config"
)

type Server struct {
	httpSrv *http.Server
}

func New(conf config.Config, handler http.Handler) *Server {
	return &Server{
		httpSrv: &http.Server{
			Addr:    conf.HTTP.Host + ":" + strconv.Itoa(conf.HTTP.Port),
			Handler: handler,
		},
	}
}

func (srv *Server) Start(ctx context.Context) error {
	go func() {
		if err := srv.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return nil
}

func (srv *Server) Stop(ctx context.Context) error {
	return srv.httpSrv.Shutdown(ctx)
}
