package http

import (
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
)

type RatingHandler struct {
	logger logging.Logger
	serv   service.Rating
}

func NewRatingHandler(logger logging.Logger, serv service.Rating) *RatingHandler {
	return &RatingHandler{
		logger: logger.With("handler", "rating", "handlerType", "http"),
		serv:   serv,
	}
}
