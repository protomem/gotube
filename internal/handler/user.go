package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/protomem/gotube/internal/access"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
	"github.com/protomem/gotube/pkg/response"
)

type User struct {
	logger logging.Logger
	serv   service.User
	accmng access.Manager
}

func NewUser(logger logging.Logger, serv service.User, accmng access.Manager) *User {
	return &User{
		logger: logger.With("handler", "user", "handlerType", "http"),
		serv:   serv,
		accmng: accmng,
	}
}

func (h *User) Get() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		const op = "handler:User.Get"

		ctx := r.Context()
		logger := h.logger.With(
			"operation", op,
			requestid.Key, requestid.Extract(ctx),
		)

		userNickname := chi.URLParam(r, "userNickname")

		user, err := h.serv.GetByNickname(ctx, userNickname)
		if err != nil {
			logger.Error("failed to get user", "error", err)

			if errors.Is(err, model.ErrUserNotFound) {
				return ErrNotFound("user")
			}

			return ErrInternal("failed to get user")
		}

		return response.Send(w, http.StatusOK, response.JSON{"user": user})
	})
}

func (h *User) Update() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{"user": "some_user"})
	})
}

func (h *User) Delete() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		const op = "handler:User.Delete"

		ctx := r.Context()
		logger := h.logger.With(
			"operation", op,
			requestid.Key, requestid.Extract(ctx),
		)

		userNickname := chi.URLParam(r, "userNickname")

		if err := h.serv.DeleteByNickname(ctx, userNickname); err != nil {
			logger.Error("failed to delete user", "error", err)

			return ErrInternal("failed to delete user")
		}

		return response.Send(w, http.StatusNoContent, nil)
	})
}

func (h *User) apiFunc(apiFn response.APIFunc) http.HandlerFunc {
	return response.BuildHandlerFunc(h.errorHandler())(apiFn)
}

func (h *User) errorHandler() response.ErrorHandler {
	return response.DefaultErrorHandler(h.logger, "handler:User.errorHandler")
}
