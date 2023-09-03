package module

import (
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/auth"
	"github.com/protomem/gotube/internal/module/common"
	"github.com/protomem/gotube/internal/module/user"
	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/logging"
)

type Modules struct {
	Common *common.Module
	User   *user.Module
	Auth   *auth.Module
}

func NewModules(logger logging.Logger, authSecret string, db *database.DB, sessmng storage.SessionManager) *Modules {
	commonMod := common.New(logger)
	userMod := user.New(logger, db)
	authMod := auth.New(logger, authSecret, sessmng, userMod.UserService)

	return &Modules{
		Common: commonMod,
		User:   userMod,
		Auth:   authMod,
	}
}
