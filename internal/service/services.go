package service

import (
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/internal/session"
)

type Services struct{}

func New(repos *repository.Repositories, sessmng session.Manager) *Services {
	return &Services{}
}
