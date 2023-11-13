package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func (handl *VideoHandler) FindNew() http.HandlerFunc {
	type Response struct {
		Videos []model.Video `json:"videos"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.VideoHandler.FindNew"
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

		opts := service.FindVideosOptions{
			Limit:  10,
			Offset: 0,
		}

		if r.URL.Query().Has("limit") {
			opts.Limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				logger.Error("failed to parse limit", "error", err)

				w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid limit",
				})

				return
			}
		}

		if r.URL.Query().Has("offset") {
			opts.Offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
			if err != nil {
				logger.Error("failed to parse offset", "error", err)

				w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid offset",
				})

				return
			}
		}

		videos, err := handl.serv.FindNew(ctx, opts)
		if err != nil {
			logger.Error("failed to find new videos", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to find new videos",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{Videos: videos})
	}
}

func (handl *VideoHandler) FindPopular() http.HandlerFunc {
	type Response struct {
		Videos []model.Video `json:"videos"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.VideoHandler.FindPopular"
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

		opts := service.FindVideosOptions{
			Limit:  10,
			Offset: 0,
		}

		if r.URL.Query().Has("limit") {
			opts.Limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				logger.Error("failed to parse limit", "error", err)

				w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid limit",
				})

				return
			}
		}

		if r.URL.Query().Has("offset") {
			opts.Offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
			if err != nil {
				logger.Error("failed to parse offset", "error", err)

				w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid offset",
				})

				return
			}
		}

		videos, err := handl.serv.FindPopular(ctx, opts)
		if err != nil {
			logger.Error("failed to find popular videos", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to find popular videos",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{Videos: videos})
	}
}

func (handl *VideoHandler) FindByAuthor() http.HandlerFunc {
	type Response struct {
		Videos []model.Video `json:"videos"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.VideoHandler.FindByAuthor"
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

		vars := mux.Vars(r)

		authorNickname, exists := vars["nickname"]
		if !exists {
			logger.Error("nickname missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "nickname missing",
			})

			return
		}

		authPayload, isAuth := jwt.Extract(ctx)

		videos, err := handl.serv.FindByAuthorNickname(ctx, authorNickname)
		if err != nil {
			logger.Error("failed to find videos by author", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to find videos by author",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		var filteredVideos []model.Video
		for _, video := range videos {
			if video.Public {
				filteredVideos = append(filteredVideos, video)
			} else {
				if isAuth && authPayload.UserID == video.Author.ID {
					filteredVideos = append(filteredVideos, video)
				}
			}
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{Videos: filteredVideos})
	}
}

func (handl *VideoHandler) FindByAuthorSubscriptions() http.HandlerFunc {
	type Response struct {
		Videos []model.Video `json:"videos"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.VideoHandler.FindByAuthorSubscriptions"
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

		vars := mux.Vars(r)

		authorNickname, exists := vars["nickname"]
		if !exists {
			logger.Error("nickname missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "nickname missing",
			})

			return
		}

		opts := service.FindVideosOptions{
			Limit:  10,
			Offset: 0,
		}

		if r.URL.Query().Has("limit") {
			opts.Limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				logger.Error("failed to parse limit", "error", err)

				w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid limit",
				})

				return
			}
		}

		if r.URL.Query().Has("offset") {
			opts.Offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
			if err != nil {
				logger.Error("failed to parse offset", "error", err)

				w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid offset",
				})
			}
		}

		videos, err := handl.serv.FindByAuthorSubscriptions(ctx, authorNickname, opts)
		if err != nil {
			logger.Error("failed to find videos by author", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to find videos by author",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{Videos: videos})
	}
}

func (handl *VideoHandler) Get() http.HandlerFunc {
	type Response struct {
		Video model.Video `json:"video"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.VideoHandler.Get"
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

		vars := mux.Vars(r)

		videoIDRaw, exists := vars["id"]
		if !exists {
			logger.Error("missing video id")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "missing video id",
			})

			return
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			logger.Error("failed to parse video id", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid video id",
			})

			return
		}

		authPayload, isAuth := jwt.Extract(ctx)

		var video model.Video
		if isAuth {
			video, err = handl.serv.Get(ctx, videoID)
		} else {
			video, err = handl.serv.GetPublic(ctx, videoID)
		}
		if err != nil {
			logger.Error("failed to get video", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to get video",
			}

			if errors.Is(err, model.ErrVideoNotFound) {
				code = http.StatusNotFound
				res["error"] = model.ErrVideoNotFound.Error()
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		if isAuth && !video.Public && authPayload.UserID != video.Author.ID {
			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusForbidden)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "access denied",
			})

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{Video: video})
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

func (handl *VideoHandler) Update() http.HandlerFunc {
	type Request struct {
		Title         *string `json:"title"`
		Description   *string `json:"description"`
		ThumbnailPath *string `json:"thumbnailPath"`
		VideoPath     *string `json:"videoPath"`
		Public        *bool   `json:"isPublic"`
	}

	type Response struct {
		Video model.Video `json:"video"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.VideoHandler.Update"
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

		vars := mux.Vars(r)

		videoIDRaw, exists := vars["id"]
		if !exists {
			logger.Error("missing video id")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "missing video id",
			})

			return
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			logger.Error("failed to parse video id", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid video id",
			})

			return
		}

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

		video, err := handl.serv.Update(ctx, videoID, service.UpdateVideoDTO(req))
		if err != nil {
			logger.Error("failed to update video", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to update video",
			}

			if errors.Is(err, model.ErrVideoNotFound) {
				code = http.StatusNotFound
				res["error"] = model.ErrVideoNotFound.Error()
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{Video: video})
	}
}

func (handl *VideoHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.VideoHandler.Delete"
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

		vars := mux.Vars(r)

		videoIDRaw, exists := vars["id"]
		if !exists {
			logger.Error("missing video id")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "missing video id",
			})

			return
		}

		videoID, err := uuid.Parse(videoIDRaw)
		if err != nil {
			logger.Error("failed to parse video id", "error", err)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid video id",
			})

			return
		}

		err = handl.serv.Delete(ctx, videoID)
		if err != nil {
			logger.Error("failed to delete video", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to delete video",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
