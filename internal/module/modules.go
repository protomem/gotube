package module

import (
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/common"
	"github.com/protomem/gotube/internal/module/user"
	"github.com/protomem/gotube/pkg/logging"
)

type Modules struct {
	Common *common.Module
	User   *user.Module
}

func NewModules(logger logging.Logger, db *database.DB) *Modules {
	return &Modules{
		Common: common.New(logger),
		User:   user.New(logger, db),
	}
}
