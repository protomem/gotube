package service

import (
	userserv "github.com/protomem/gotube/internal/module/user/service"
	"github.com/protomem/gotube/internal/storage"
)

var _ AuthService = (*AuthServiceImpl)(nil)

type (
	AuthService interface{}

	AuthServiceImpl struct {
		secret string

		sessmng storage.SessionManager

		userServ userserv.UserService
	}
)

func NewAuthService(secret string, sessmng storage.SessionManager, userServ userserv.UserService) *AuthServiceImpl {
	return &AuthServiceImpl{
		secret:   secret,
		sessmng:  sessmng,
		userServ: userServ,
	}
}
