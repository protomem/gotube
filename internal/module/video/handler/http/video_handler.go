package http

import (
	"github.com/protomem/gotube/internal/module/video/service"
	"github.com/protomem/gotube/pkg/logging"
)

type VideoHandler struct {
	logger logging.Logger

	videoServ service.VideoService
}

func NewVideoHandler(logger logging.Logger, videoServ service.VideoService) *VideoHandler {
	return &VideoHandler{
		logger:    logger.With("handler", "video", "handlerType", "http"),
		videoServ: videoServ,
	}
}
