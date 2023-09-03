package common

import (
	handlhttp "github.com/protomem/gotube/internal/module/common/handler/http"
	mdwhttp "github.com/protomem/gotube/internal/module/common/middleware/http"
	"github.com/protomem/gotube/pkg/logging"
)

type Module struct {
	*handlhttp.CommonHandler
	*mdwhttp.CommonMiddleware
}

func New(logger logging.Logger) *Module {
	logger = logger.With("module", "common")

	return &Module{
		CommonHandler:    handlhttp.NewCommonHandler(logger),
		CommonMiddleware: mdwhttp.NewCommonMiddleware(logger),
	}
}
