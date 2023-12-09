package service

import (
	"github.com/protomem/gotube/internal/session"
)

type (
	Auth interface{}

	AuthImpl struct {
		userServ User
		sessmng  session.Manager
	}
)

func NewAuth(userServ User, sessmng session.Manager) *AuthImpl {
	return &AuthImpl{
		userServ: userServ,
		sessmng:  sessmng,
	}
}
