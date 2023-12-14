package handler

import (
	"encoding/json"
	"net/http"

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

type Comment struct {
	logger logging.Logger
	serv   service.Comment
	accmng access.Manager
}

func NewComment(logger logging.Logger, serv service.Comment, accmng access.Manager) *Comment {
	return &Comment{
		logger: logger.With("handler", "comment", "handlerType", "http"),
		serv:   serv,
		accmng: accmng,
	}
}

func (h *Comment) List() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		const op = "handler:Comment.List"

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

		videos, err := h.serv.FindByVideoID(ctx, videoID)
		if err != nil {
			logger.Error("failed to find comments", "error", err)

			return ErrInternal("failed to find comments")
		}

		return response.Send(w, http.StatusOK, response.JSON{
			"videos": videos,
		})
	})
}

func (h *Comment) Create() http.HandlerFunc {
	type Request struct {
		Message string `json:"message"`
	}

	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		const op = "handler:Comment.Create"

		ctx := r.Context()
		logger := h.logger.With(
			"operation", op,
			requestid.Key, requestid.Extract(ctx),
		)

		authPayload, _ := jwt.Extract(ctx)

		videoID, err := h.extractVideoIDFromRequest(r)
		if err != nil {
			logger.Error("failed to extract video id", "error", err)

			return ErrBadRequest
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode request", "error", err)

			return ErrBadRequest
		}

		comment, err := h.serv.Create(ctx, service.CreateCommentDTO{
			Message:  req.Message,
			AuthorID: authPayload.UserID,
			VideoID:  videoID,
		})
		if err != nil {
			logger.Error("failed to create comment", "error", err)

			return ErrInternal("failed to create comment")
		}

		return response.Send(w, http.StatusOK, response.JSON{
			"comment": comment,
		})
	})
}

func (h *Comment) Update() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{
			"comment": "some_comment",
		})
	})
}

func (h *Comment) Delete() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		const op = "handler:Comment.Delete"

		ctx := r.Context()
		logger := h.logger.With(
			"operation", op,
			requestid.Key, requestid.Extract(ctx),
		)

		commentID, err := h.extractCommentIDFromRequest(r)
		if err != nil {
			logger.Error("failed to extract comment id", "error", err)

			return ErrBadRequest
		}

		if err := h.serv.Delete(ctx, commentID); err != nil {
			logger.Error("failed to delete comment", "error", err)

			return ErrInternal("failed to delete comment")
		}

		return response.Send(w, http.StatusNoContent, nil)
	})
}

func (h *Comment) apiFunc(apiFn response.APIFunc) http.HandlerFunc {
	return response.BuildHandlerFunc(h.errorHandler())(apiFn)
}

func (h *Comment) errorHandler() response.ErrorHandler {
	return response.DefaultErrorHandler(h.logger, "handler:Auth.errorHandler")
}

func (h *Comment) extractVideoIDFromRequest(r *http.Request) (model.ID, error) {
	videoIDRaw := chi.URLParam(r, "videoId")

	videoID, err := uuid.Parse(videoIDRaw)
	if err != nil {
		return model.ID{}, err
	}

	return videoID, nil
}

func (h *Comment) extractCommentIDFromRequest(r *http.Request) (model.ID, error) {
	commentIDRaw := chi.URLParam(r, "commentId")

	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return model.ID{}, err
	}

	return commentID, nil
}
