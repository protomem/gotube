package middleware

import (
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
)

type Middlewares struct {
	*Common
	*Auth
}

func New(logger logging.Logger, servs *service.Services) *Middlewares {
	return &Middlewares{
		Common: NewCommon(),
		Auth:   NewAuth(logger, servs.Auth),
	}
}
