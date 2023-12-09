package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/protomem/gotube/internal/access"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
	"github.com/protomem/gotube/pkg/response"
)

type Auth struct {
	logger logging.Logger
	serv   service.Auth
	accmng access.Manager
}

func NewAuth(logger logging.Logger, serv service.Auth, accmng access.Manager) *Auth {
	return &Auth{
		logger: logger.With("handler", "auth", "handlerType", "http"),
		serv:   serv,
		accmng: accmng,
	}
}

func (h *Auth) Register() http.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		const op = "handler:Auth.Register"

		ctx := r.Context()
		logger := h.logger.With(
			"operation", op,
			requestid.Key, requestid.Extract(ctx),
		)

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode request body", "error", err)

			return ErrBadRequest
		}

		user, token, err := h.serv.Register(ctx, service.RegisterDTO(req))
		if err != nil {
			logger.Error("failed to register", "error", err)

			if errors.Is(err, model.ErrUserAlreadyExists) {
				return ErrConflict("user")
			}

			return ErrInternal("failed to register")
		}

		return response.Send(w, http.StatusCreated, response.JSON{
			"accessToken":  token.Access,
			"refreshToken": token.Refresh,
			"user":         user,
		})
	})
}

func (h *Auth) Login() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		const op = "handler:Auth.Login"

		ctx := r.Context()
		logger := h.logger.With(
			"operation", op,
			requestid.Key, requestid.Extract(ctx),
		)

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode request body", "error", err)

			return ErrBadRequest
		}

		user, token, err := h.serv.Login(ctx, service.LoginDTO(req))
		if err != nil {
			logger.Error("failed to login", "error", err)

			if errors.Is(err, model.ErrUserNotFound) {
				return ErrNotFound("user")
			}

			return ErrInternal("failed to login")
		}

		return response.Send(w, http.StatusOK, response.JSON{
			"accessToken":  token.Access,
			"refreshToken": token.Refresh,
			"user":         user,
		})
	})
}

func (h *Auth) RefreshToken() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{"user": "some_user"})
	})
}

func (h *Auth) Logout() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusNoContent, nil)
	})
}

func (h *Auth) apiFunc(apiFn response.APIFunc) http.HandlerFunc {
	return response.BuildHandlerFunc(h.errorHandler())(apiFn)
}

func (h *Auth) errorHandler() response.ErrorHandler {
	return response.DefaultErrorHandler(h.logger, "handler:Auth.errorHandler")
}
