package http

import (
	"github.com/protomem/gotube/internal/module/media/service"
	"github.com/protomem/gotube/pkg/logging"
)

type MediaHandler struct {
	logger logging.Logger

	mediaServ service.MediaService
}

func NewMediaHandler(logger logging.Logger, mediaServ service.MediaService) *MediaHandler {
	return &MediaHandler{
		logger:    logger.With("handler", "media", "handlerType", "http"),
		mediaServ: mediaServ,
	}
}
