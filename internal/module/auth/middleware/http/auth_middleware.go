package http

import (
	"github.com/protomem/gotube/internal/module/auth/service"
	"github.com/protomem/gotube/pkg/logging"
)

type AuthMiddleware struct {
	logger logging.Logger

	authServ service.AuthService
}

func NewAuthMiddleware(logger logging.Logger, authServ service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		logger:   logger.With("middleware", "auth", "middlewareType", "http"),
		authServ: authServ,
	}
}
