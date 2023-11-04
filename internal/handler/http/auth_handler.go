package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httpheader"
	"github.com/protomem/gotube/pkg/logging"
)

type AuthHandler struct {
	logger logging.Logger
	serv   service.Auth
}

func NewAuthHandler(logger logging.Logger, serv service.Auth) *AuthHandler {
	return &AuthHandler{
		logger: logger.With("handler", "auth", "handlerType", "http"),
		serv:   serv,
	}
}

func (handl *AuthHandler) Register() http.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type Response struct {
		model.PairTokens
		User model.User `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.AuthHandler.Register"
		var err error

		ctx := r.Context()
		logger := handl.logger.With("operation", op)

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

		user, tokens, err := handl.serv.Register(ctx, service.RegisterDTO(req))
		if err != nil {
			logger.Error("failed to create user", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to create user",
			}

			if errors.Is(err, model.ErrUserExists) {
				code = http.StatusConflict
				res = map[string]string{
					"error": model.ErrUserExists.Error(),
				}
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Response{
			PairTokens: tokens,
			User:       user,
		})
	}
}

func (handl *AuthHandler) Login() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		model.PairTokens
		User model.User `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.AuthHandler.Login"
		var err error

		ctx := r.Context()
		logger := handl.logger.With("operation", op)

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

		user, tokens, err := handl.serv.Login(ctx, service.LoginDTO(req))
		if err != nil {
			logger.Error("failed to login user", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to login user",
			}

			if errors.Is(err, model.ErrUserNotFound) {
				code = http.StatusNotFound
				res = map[string]string{
					"error": model.ErrUserNotFound.Error(),
				}
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{
			PairTokens: tokens,
			User:       user,
		})
	}
}

func (handl *AuthHandler) Refresh() http.HandlerFunc {
	type Request struct {
		RefreshToken string `json:"refreshToken"`
	}

	type Response struct {
		model.PairTokens
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.AuthHandler.Refresh"
		var err error

		ctx := r.Context()
		logger := handl.logger.With("operation", op)

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

		tokens, err := handl.serv.RefreshTokens(ctx, req.RefreshToken)
		if err != nil {
			logger.Error("failed to refresh tokens", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to refresh tokens",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{
			PairTokens: tokens,
		})
	}
}

func (handl *AuthHandler) Logout() http.HandlerFunc {
	type Request struct {
		RefreshToken string `json:"refreshToken"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.AuthHandler.Logout"
		var err error

		ctx := r.Context()
		logger := handl.logger.With("operation", op)

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

		err = handl.serv.Logout(ctx, req.RefreshToken)
		if err != nil {
			logger.Error("failed to logout user", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to logout user",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
