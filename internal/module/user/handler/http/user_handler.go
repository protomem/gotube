package http

import (
	"github.com/protomem/gotube/internal/module/user/service"
	"github.com/protomem/gotube/pkg/logging"
)

type UserHandler struct {
	logger logging.Logger

	userServ service.UserService
}

func NewUserHandler(logger logging.Logger, userServ service.UserService) *UserHandler {
	return &UserHandler{
		logger:   logger.With("handler", "user", "handlerType", "http"),
		userServ: userServ,
	}
}
