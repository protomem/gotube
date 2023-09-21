package http

import (
	"github.com/protomem/gotube/internal/module/video/service"
	"github.com/protomem/gotube/pkg/logging"
)

type RatingHandler struct {
	logger logging.Logger

	ratingServ service.RatingService
}

func NewRatingHandler(logger logging.Logger, ratingServ service.RatingService) *RatingHandler {
	return &RatingHandler{
		logger:     logger.With("handler", "rating", "handlerType", "http"),
		ratingServ: ratingServ,
	}
}
