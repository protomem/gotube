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

func (h *VideoHandler) HandleGetVideo() echo.HandlerFunc {
	type Request struct {
		ID uuid.UUID `param:"videoId"`
	}

	type Response struct {
		Video model.Video `json:"video"`
	}

	return func(c echo.Context) error {
		const op = "VideoHandler.HandleGetVideo"
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

		video, err := h.videoServ.FindOneVideo(ctx, req.ID)
		if err != nil {
			logger.Error("failed to find video", "error", err)

			if errors.Is(err, model.ErrVideoNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, model.ErrVideoNotFound)
			}

			return echo.ErrInternalServerError
		}

		if !video.Public && video.User.ID != authPayload.UserID {
			return echo.NewHTTPError(http.StatusForbidden, "access denied")
		}

		return c.JSON(http.StatusOK, Response{Video: video})
	}
}

func (h *VideoHandler) HandleGetAllVideos() echo.HandlerFunc {
	type Request struct {
		Limit         uint64 `query:"limit"`
		Offset        uint64 `query:"offset"`
		SortBy        string `query:"sort_by"`
		Subscriptions bool   `query:"is_subs"`
		UserNickname  string `query:"user"`
		Query         string `query:"query"`
	}

	type Response struct {
		Videos []model.Video `json:"videos"`
	}

	return func(c echo.Context) error {
		const op = "VideoHandler.HandleGetAllVideos"
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

		authPayload, authPayloadOk := jwt.Extract(ctx)

		if req.Limit == 0 {
			req.Limit = 10
		}

		if req.SortBy != "new" && req.SortBy != "popular" {
			req.SortBy = "new"
		}

		opts := service.FindAllVideosOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		}

		var videos []model.Video
		if req.Query != "" {
			videos, err = h.videoServ.SearchVideos(ctx, req.Query, opts)
		} else if req.UserNickname != "" {
			videos, err = h.videoServ.FindAllNewVideosByUserNickname(ctx, req.UserNickname)
		} else if req.Subscriptions {
			if !authPayloadOk {
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}

			videos, err = h.videoServ.FindAllPublicNewVideosFromSubscriptions(ctx, authPayload.UserID, opts)
		} else {
			if req.SortBy == "new" {
				videos, err = h.videoServ.FindAllPublicNewVideos(ctx, opts)
			} else {
				videos, err = h.videoServ.FindAllPublicPopularVideos(ctx, opts)
			}
		}
		if err != nil {
			logger.Error("failed to find videos", "error", err)

			return echo.ErrInternalServerError
		}

		// Auth filter
		filteredVideos := make([]model.Video, 0, len(videos))
		for _, video := range videos {
			if video.Public ||
				(authPayloadOk && !video.Public && video.User.ID == authPayload.UserID) {
				filteredVideos = append(filteredVideos, video)
			}
		}

		return c.JSON(http.StatusOK, Response{Videos: filteredVideos})
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
