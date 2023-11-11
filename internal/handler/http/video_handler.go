package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httpheader"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
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

func (handl *VideoHandler) Create() http.HandlerFunc {
	type Request struct {
		Title         string  `json:"title"`
		Description   *string `json:"description"`
		ThumbnailPath string  `json:"thumbnailPath"`
		VideoPath     string  `json:"videoPath"`
		Public        *bool   `json:"isPublic"`
	}

	type Response struct {
		Video model.Video `json:"video"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.VideoHandler.Create"
		var err error

		ctx := r.Context()
		logger := handl.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		var req Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("failed to decode request", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to decode request",
			})

			return
		}

		authPayload, _ := jwt.Extract(ctx)

		video, err := handl.serv.Create(ctx, service.CreateVideoDTO{
			Title:         req.Title,
			Description:   req.Description,
			ThumbnailPath: req.ThumbnailPath,
			VideoPath:     req.VideoPath,
			Public:        req.Public,
			AuthorID:      authPayload.UserID,
		})
		if err != nil {
			logger.Error("failed to create video", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to create video",
			}

			if errors.Is(err, model.ErrVideoExists) {
				code = http.StatusConflict
				res["error"] = model.ErrVideoExists.Error()
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Response{Video: video})
	}
}
