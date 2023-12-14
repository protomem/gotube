package handler

import (
	"github.com/protomem/gotube/internal/access"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/logging"
)

type Handlers struct {
	*Common

	*User
	*Auth
	*Video
	*Comment
}

func New(logger logging.Logger, servs *service.Services, store storage.Storage, accmng access.Manager) *Handlers {
	return &Handlers{
		Common: NewCommon(logger),

		User:    NewUser(logger, servs.User, accmng),
		Auth:    NewAuth(logger, servs.Auth, accmng),
		Video:   NewVideo(logger, servs.Video, accmng),
		Comment: NewComment(logger, servs.Comment, accmng),
	}
}
