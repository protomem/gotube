package service

import (
	"github.com/protomem/gotube/internal/config"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/hashing"
)

type Services struct {
	User
	Auth
	Subscription
	Video
}

func New(authConf config.Auth, repos *repository.Repositories, hasher hashing.Hasher) *Services {
	var (
		user  = NewUser(repos.User, hasher)
		auth  = NewAuth(authConf, user)
		sub   = NewSubscription(repos.Subscription, user)
		video = NewVideo(repos.Video, user)
	)

	return &Services{
		User:         user,
		Auth:         auth,
		Subscription: sub,
		Video:        video,
	}
}
