package service

import (
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/hashing"
)

type Services struct {
	User
}

func New(repos *repository.Repositories, hasher hashing.Hasher) *Services {
	return &Services{
		User: NewUser(repos.User, hasher),
	}
}
