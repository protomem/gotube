package middleware

import (
	"net/http"

	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
)

type MiddlewareFunc func(next http.Handler) http.Handler

type Middlewares struct {
	*Common
	*Auth
}

func New(logger logging.Logger, servs *service.Services) *Middlewares {
	return &Middlewares{
		Common: NewCommon(logger),
		Auth:   NewAuth(logger, servs.Auth),
	}
}
