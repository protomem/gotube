package module

import (
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/auth"
	"github.com/protomem/gotube/internal/module/common"
	"github.com/protomem/gotube/internal/module/media"
	"github.com/protomem/gotube/internal/module/subscription"
	"github.com/protomem/gotube/internal/module/user"
	"github.com/protomem/gotube/internal/module/video"
	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/logging"
)

type Modules struct {
	Common       *common.Module
	User         *user.Module
	Auth         *auth.Module
	Subscription *subscription.Module
	Video        *video.Module

	Media *media.Module
}

func NewModules(
	logger logging.Logger,
	authSecret string,
	db *database.DB,
	bstore storage.BlobStorage,
	sessmng storage.SessionManager,
) *Modules {
	commonMod := common.New(logger)

	userMod := user.New(logger, db)
	authMod := auth.New(logger, authSecret, sessmng, userMod.UserService)
	subMod := subscription.New(logger, db, userMod.UserService)
	videoMod := video.New(logger, db)

	mediaMod := media.New(logger, bstore)

	return &Modules{
		Common:       commonMod,
		User:         userMod,
		Auth:         authMod,
		Subscription: subMod,
		Video:        videoMod,
		Media:        mediaMod,
	}
}
