package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/protomem/gotube/internal/access"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
	"github.com/protomem/gotube/pkg/response"
)

type Video struct {
	logger logging.Logger
	serv   service.Video
	accmng access.Manager
}

func NewVideo(logger logging.Logger, serv service.Video, accmng access.Manager) *Video {
	return &Video{
		logger: logger.With("handler", "video", "handlerType", "http"),
		serv:   serv,
		accmng: accmng,
	}
}

func (h *Video) List() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{"videos": "some_videos"})
	})
}

func (h *Video) Get() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{"video": "some_video"})
	})
}

func (h *Video) Create() http.HandlerFunc {
	type Request struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		ThumbnailPath string `json:"thumbnailPath"`
		VideoPath     string `json:"videoPath"`
		Public        bool   `json:"isPublic"`
	}

	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		const op = "handler:Video.Create"

		ctx := r.Context()
		logger := h.logger.With(
			"operation", op,
			requestid.Key, requestid.Extract(ctx),
		)

		authPayload, _ := jwt.Extract(ctx)

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode request", "error", err)

			return ErrBadRequest
		}

		video, err := h.serv.Create(ctx, service.CreateVideoDTO{
			Title:         req.Title,
			Description:   req.Description,
			ThumbnailPath: req.ThumbnailPath,
			VideoPath:     req.VideoPath,
			AuthorID:      authPayload.UserID,
			Public:        req.Public,
		})
		if err != nil {
			logger.Error("failed to create video", "error", err)

			if errors.Is(err, model.ErrVideoAlreadyExists) {
				return ErrConflict("video")
			}

			return ErrInternal("failed to create video")
		}

		return response.Send(w, http.StatusCreated, response.JSON{
			"video": video,
		})
	})
}

func (h *Video) Update() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{"video": "some_video"})
	})
}

func (h *Video) Delete() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusNoContent, nil)
	})
}

func (h *Video) apiFunc(apiFn response.APIFunc) http.HandlerFunc {
	return response.BuildHandlerFunc(h.errorHandler())(apiFn)
}

func (h *Video) errorHandler() response.ErrorHandler {
	return response.DefaultErrorHandler(h.logger, "handler:User.errorHandler")
}
