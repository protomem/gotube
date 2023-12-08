package service

import (
	"github.com/protomem/gotube/internal/hashing"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/internal/session"
)

type Services struct {
	User
}

func New(repos *repository.Repositories, sessmng session.Manager) *Services {
	return &Services{
		User: NewUser(repos.User, hashing.NewBcrypt(hashing.BcryptDefaultCost)),
	}
}
