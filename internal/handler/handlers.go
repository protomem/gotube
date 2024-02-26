package handler

import (
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
)

type Handlers struct {
	*Common
	*User
	*Auth
	*Subscription
	*Video
	*Rating
	*Comment
}

func New(logger logging.Logger, servs *service.Services) *Handlers {
	return &Handlers{
		Common:       NewCommon(),
		User:         NewUser(logger, servs.User),
		Auth:         NewAuth(logger, servs.Auth),
		Subscription: NewSubscription(logger, servs.Subscription),
		Video:        NewVideo(logger, servs.Video),
		Rating:       NewRating(logger, servs.Rating),
		Comment:      NewComment(logger, servs.Comment),
	}
}
