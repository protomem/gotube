package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httpheader"
	"github.com/protomem/gotube/pkg/logging"
)

type UserHandler struct {
	logger logging.Logger
	serv   service.User
}

func NewUserHandler(logger logging.Logger, serv service.User) *UserHandler {
	return &UserHandler{
		logger: logger.With("handler", "user", "handlerType", "http"),
		serv:   serv,
	}
}

func (handl *UserHandler) Get() http.HandlerFunc {
	type Response struct {
		User model.User `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.UserHandler.Get"
		var err error

		ctx := r.Context()
		logger := handl.logger.With("operation", op)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		vars := mux.Vars(r)

		userNickname, exists := vars["nickname"]
		if !exists {
			logger.Error("nickname missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "nickname missing",
			})

			return
		}

		user, err := handl.serv.GetByNickname(ctx, userNickname)
		if err != nil {
			logger.Error("failed to get user", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to get user",
			}

			if errors.Is(err, model.ErrUserNotFound) {
				code = http.StatusNotFound
				res = map[string]string{
					"error": model.ErrUserNotFound.Error(),
				}
			}

			logger.Debug("test", "res", res)

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{User: user})
	}
}

func (handl *UserHandler) Create() http.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type Response struct {
		User model.User `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.UserHandler.Create"
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

		user, err := handl.serv.Create(ctx, service.CreateUserDTO(req))
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
		err = json.NewEncoder(w).Encode(Response{User: user})
	}
}

func (handl *UserHandler) Update() http.HandlerFunc {
	type Request struct {
		Nickname    *string `json:"nickname"`
		Email       *string `json:"email"`
		AvatarPath  *string `json:"avatarPath"`
		Description *string `json:"description"`

		OldPassword *string `json:"oldPassword"`
		NewPassword *string `json:"newPassword"`
	}

	type Response struct {
		User model.User `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.UserHandler.Update"
		var err error

		ctx := r.Context()
		logger := handl.logger.With("operation", op)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		vars := mux.Vars(r)

		userNickname, exists := vars["nickname"]
		if !exists {
			logger.Error("nickname missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "nickname missing",
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

		user, err := handl.serv.UpdateByNickname(ctx, userNickname, service.UpdateUserDTO(req))
		if err != nil {
			logger.Error("failed to update user", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to update user",
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
		err = json.NewEncoder(w).Encode(Response{User: user})
	}
}

func (handl *UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.UserHandler.Delete"
		var err error

		ctx := r.Context()
		logger := handl.logger.With("operation", op)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		vars := mux.Vars(r)

		userNickname, exists := vars["nickname"]
		if !exists {
			logger.Error("nickname missing")

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "nickname missing",
			})

			return
		}

		err = handl.serv.DeleteByNickname(ctx, userNickname)
		if err != nil {
			logger.Error("failed to delete user", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to delete user",
			}

			w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
