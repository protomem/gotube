package http

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/module/video/model"
	"github.com/protomem/gotube/internal/module/video/service"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
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

func (h *RatingHandler) HandleLike() echo.HandlerFunc {
	type Request struct {
		VideoID uuid.UUID `param:"videoId"`
	}

	return func(c echo.Context) error {
		const op = "RatingHandler.HandleLike"
		var err error

		ctx := c.Request().Context()
		logger := h.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		var req Request
		err = c.Bind(&req)
		if err != nil {
			logger.Error("failed to bind request", "error", err)

			return echo.ErrBadRequest
		}

		authPaylod, _ := jwt.Extract(ctx)

		err = h.ratingServ.Like(ctx, service.LikeDTO{
			UserID:  authPaylod.UserID,
			VideoID: req.VideoID,
		})
		if err != nil {
			logger.Error("failed to like video", "error", err)

			if errors.Is(err, model.ErrRatingAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, model.ErrRatingAlreadyExists.Error())
			}

			return echo.ErrInternalServerError
		}

		return c.NoContent(http.StatusCreated)
	}
}

func (h *RatingHandler) HandleDislike() echo.HandlerFunc {
	type Request struct {
		VideoID uuid.UUID `param:"videoId"`
	}

	return func(c echo.Context) error {
		const op = "RatingHandler.HandleDislike"
		var err error

		ctx := c.Request().Context()
		logger := h.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		var req Request
		err = c.Bind(&req)
		if err != nil {
			logger.Error("failed to bind request", "error", err)

			return echo.ErrBadRequest
		}

		authPaylod, _ := jwt.Extract(ctx)

		err = h.ratingServ.Dislike(ctx, service.DislikeDTO{
			UserID:  authPaylod.UserID,
			VideoID: req.VideoID,
		})
		if err != nil {
			logger.Error("failed to dislike video", "error", err)

			if errors.Is(err, model.ErrRatingAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, model.ErrRatingAlreadyExists.Error())
			}

			return echo.ErrInternalServerError
		}

		return c.NoContent(http.StatusCreated)
	}
}
