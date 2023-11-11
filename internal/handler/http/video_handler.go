package http

import (
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
)

type VideoHandler struct {
	logger logging.Logger
	serv   service.Video
}

func NewVideoHandler(logger logging.Logger, serv service.Video) *VideoHandler {
	return &VideoHandler{
		logger: logger.With("handler", "video", "handlerType", "http"),
		serv:   serv,
	}
}
