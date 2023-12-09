package service

import (
	"github.com/protomem/gotube/internal/hashing"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/internal/session"
)

type Services struct {
	User
	Auth
}

func New(authSecret string, repos *repository.Repositories, sessmng session.Manager) *Services {
	userServ := NewUser(repos.User, hashing.NewBcrypt(hashing.BcryptDefaultCost))
	authServ := NewAuth(authSecret, userServ, sessmng)

	return &Services{
		User: userServ,
		Auth: authServ,
	}
}
