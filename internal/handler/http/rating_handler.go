package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httpheader"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
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

func (handl *RatingHandler) Like() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.RatingHandler.Like"
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

		authPayload, _ := jwt.Extract(ctx)

		vars := mux.Vars(r)

		videoIDRaw, exists := vars["id"]
		if !exists {
			logger.Error("video id missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "video id missing",
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

		err = handl.serv.Like(ctx, service.LikeDTO{
			VideoID: videoID,
			UserID:  authPayload.UserID,
		})
		if err != nil {
			logger.Error("failed to like video", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to like video",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (handl *RatingHandler) Dislike() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.RatingHandler.Dislike"
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

		authPayload, _ := jwt.Extract(ctx)

		vars := mux.Vars(r)

		videoIDRaw, exists := vars["id"]
		if !exists {
			logger.Error("video id missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "video id missing",
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

		err = handl.serv.Dislike(ctx, service.DislikeDTO{
			VideoID: videoID,
			UserID:  authPayload.UserID,
		})
		if err != nil {
			logger.Error("failed to dislike video", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to dislike video",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (handl *RatingHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.RatingHandler.Delete"
		var err error

		ctx := r.Context()
		logger := handl.logger.With(
			"operation", op,
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		authPayload, _ := jwt.Extract(ctx)

		vars := mux.Vars(r)

		videoIDRaw, exists := vars["id"]
		if !exists {
			logger.Error("video id missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "video id missing",
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

		err = handl.serv.DeleteByVideoIDAndUserID(ctx, videoID, authPayload.UserID)
		if err != nil {
			logger.Error("failed to delete rating", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to delete rating",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
