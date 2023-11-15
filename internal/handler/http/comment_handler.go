package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httpheader"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type CommentHandler struct {
	logger logging.Logger
	serv   service.Comment
}

func NewCommentHandler(logger logging.Logger, serv service.Comment) *CommentHandler {
	return &CommentHandler{
		logger: logger.With("handler", "comment", "handlerType", "http"),
		serv:   serv,
	}
}

func (handl *CommentHandler) FindByVideoID() http.HandlerFunc {
	type Response struct {
		Comments []model.Comment `json:"comments"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.CommentHandler.FindByVideoID"
		var err error

		ctx := r.Context()
		logger := handl.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "err", err)
			}
		}()

		videoIDRaw, exists := mux.Vars(r)["id"]
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

		comments, err := handl.serv.FindByVideoID(ctx, videoID)
		if err != nil {
			logger.Error("failed to find comments", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to find comments",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{Comments: comments})
	}
}

func (handl *CommentHandler) Create() http.HandlerFunc {
	type Request struct {
		Content string `json:"content"`
	}

	type Response struct {
		Comment model.Comment `json:"comment"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.CommentHandler.Create"
		var err error

		ctx := r.Context()
		logger := handl.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "err", err)
			}
		}()

		authPayload, _ := jwt.Extract(ctx)

		videoIDRaw, exists := mux.Vars(r)["id"]
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

		comment, err := handl.serv.Create(ctx, service.CreateCommentDTO{
			Content:  req.Content,
			AuthorID: authPayload.UserID,
			VideoID:  videoID,
		})
		if err != nil {
			logger.Error("failed to create comment", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"err": "failed to create comment",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Response{Comment: comment})
	}
}
