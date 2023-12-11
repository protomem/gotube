package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func (h *Video) ListNew() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		const op = "handler:Video.ListNew"

		ctx := r.Context()
		logger := h.logger.With(
			"operation", op,
			requestid.Key, requestid.Extract(ctx),
		)

		opts, err := h.extractLimitAndOffsetFromRequest(r)
		if err != nil {
			logger.Error("failed to extract limit and offset", "error", err)

			return ErrBadRequest
		}

		videos, err := h.serv.FindNew(ctx, opts)
		if err != nil {
			logger.Error("failed to find new videos", "error", err)

			return ErrInternal("failed to find new videos")
		}

		return response.Send(w, http.StatusOK, response.JSON{
			"videos": videos,
		})
	})
}

func (h *Video) Get() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		const op = "handler:Video.Get"

		ctx := r.Context()
		logger := h.logger.With(
			"operation", op,
			requestid.Key, requestid.Extract(ctx),
		)

		videoID, err := h.extractVideoIDFromRequest(r)
		if err != nil {
			logger.Error("failed to extract video id", "error", err)

			return ErrBadRequest
		}

		video, err := h.serv.Get(ctx, videoID)
		if err != nil {
			logger.Error("failed to get video", "error", err)

			if errors.Is(err, model.ErrVideoNotFound) {
				return ErrNotFound("video")
			}

			return ErrInternal("failed to get video")
		}

		return response.Send(w, http.StatusOK, response.JSON{
			"video": video,
		})
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

func (h *Video) extractVideoIDFromRequest(r *http.Request) (model.ID, error) {
	videoIDRaw := chi.URLParam(r, "videoId")

	videoID, err := uuid.Parse(videoIDRaw)
	if err != nil {
		return model.ID{}, err
	}

	return videoID, nil
}

func (h *Video) extractLimitAndOffsetFromRequest(r *http.Request) (service.FindOptions, error) {
	opts := service.FindOptions{
		Limit:  10,
		Offset: 0,
	}

	if r.URL.Query().Has("limit") {
		var err error
		opts.Limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			return service.FindOptions{}, err
		}
	}

	if r.URL.Query().Has("offset") {
		var err error
		opts.Offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
		if err != nil {
			return service.FindOptions{}, err
		}
	}

	return opts, nil
}
