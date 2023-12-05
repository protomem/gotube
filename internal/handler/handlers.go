package handler

import "github.com/protomem/gotube/pkg/logging"

type Handlers struct {
	*Common
}

func New(logger logging.Logger) *Handlers {
	return &Handlers{
		Common: NewCommon(logger),
	}
}
