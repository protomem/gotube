package handler

import (
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
)

type Handlers struct {
	*Common
	*User
	*Auth
}

func New(logger logging.Logger, servs *service.Services) *Handlers {
	return &Handlers{
		Common: NewCommon(),
		User:   NewUser(logger, servs.User),
		Auth:   NewAuth(logger, servs.Auth),
	}
}
