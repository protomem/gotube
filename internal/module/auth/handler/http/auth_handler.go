package http

import (
	"github.com/protomem/gotube/internal/module/auth/service"
	"github.com/protomem/gotube/pkg/logging"
)

type AuthHandler struct {
	logger logging.Logger

	authServ service.AuthService
}

func NewAuthHandler(logger logging.Logger, authServ service.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:   logger.With("handler", "auth", "handlerType", "http"),
		authServ: authServ,
	}
}
