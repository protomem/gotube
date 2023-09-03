package module

import (
	"github.com/protomem/gotube/internal/module/common"
	"github.com/protomem/gotube/pkg/logging"
)

type Modules struct {
	Common *common.Module
}

func NewModules(logger logging.Logger) *Modules {
	return &Modules{
		Common: common.New(logger),
	}
}
