package service

import (
	"github.com/protomem/gotube/internal/hashing"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/internal/session"
)

type Services struct {
	User
	Auth
	Video
	Comment
}

func New(authSecret string, repos *repository.Repositories, sessmng session.Manager) *Services {
	userServ := NewUser(repos.User, hashing.NewBcrypt(hashing.BcryptDefaultCost))
	authServ := NewAuth(authSecret, userServ, sessmng)
	videoServ := NewVideo(repos.Video)
	commentServ := NewComment(repos.Comment, userServ, videoServ)

	return &Services{
		User:    userServ,
		Auth:    authServ,
		Video:   videoServ,
		Comment: commentServ,
	}
}
