package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/module/video/model"
	"github.com/protomem/gotube/internal/module/video/service"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
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

func (h *VideoHandler) HandleCreateVideo() echo.HandlerFunc {
	type Request struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		ThumbnailPath string `json:"thumbnailPath"`
		VideoPath     string `json:"videoPath"`
	}

	type Response struct {
		Video model.Video `json:"video"`
	}

	return func(c echo.Context) error {
		const op = "VideoHandler.HandleCreateVideo"
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

		authPayload, _ := jwt.Extract(ctx)

		video, err := h.videoServ.CreateVideo(ctx, service.CreateVideoDTO{
			Title:         req.Title,
			Description:   req.Description,
			ThumbnailPath: req.ThumbnailPath,
			VideoPath:     req.VideoPath,
			UserID:        authPayload.UserID,
		})
		if err != nil {
			logger.Error("failed to create video", "error", err)

			if errors.Is(err, model.ErrVideoAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, model.ErrVideoAlreadyExists)
			}

			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusCreated, Response{Video: video})
	}
}
