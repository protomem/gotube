package auth

import (
	handlhttp "github.com/protomem/gotube/internal/module/auth/handler/http"
	mdwhttp "github.com/protomem/gotube/internal/module/auth/middleware/http"
	"github.com/protomem/gotube/internal/module/auth/service"
	userserv "github.com/protomem/gotube/internal/module/user/service"
	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/logging"
)

type Module struct {
	*handlhttp.AuthHandler
	*mdwhttp.AuthMiddleware
	service.AuthService
}

func New(logger logging.Logger, secret string, sessmng storage.SessionManager, userServ userserv.UserService) *Module {
	logger = logger.With("module", "auth")

	authServ := service.NewAuthService(secret, sessmng, userServ)
	authMdw := mdwhttp.NewAuthMiddleware(logger, authServ)
	authHdl := handlhttp.NewAuthHandler(logger, authServ)

	return &Module{
		AuthHandler:    authHdl,
		AuthMiddleware: authMdw,
		AuthService:    authServ,
	}
}
