package handler

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

type User struct {
	logger logging.Logger
	serv   service.User
}

func NewUser(logger logging.Logger, serv service.User) *User {
	return &User{
		logger: logger.With("handler", "user"),
		serv:   serv,
	}
}

func (h *User) Get() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		userNickname, ok := mux.Vars(r)["userNickname"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing nickname")
		}

		user, err := h.serv.GetByNickname(r.Context(), userNickname)
		if err != nil {
			return err
		}

		return httplib.WriteJSON(w, http.StatusOK, httplib.JSON{"user": user})
	}, h.errorHandler("handler.User.Get"))
}

func (h *User) Create() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		var request struct {
			Nickname string `json:"nickname"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := httplib.DecodeJSON(r, &request); err != nil {
			return err
		}

		user, err := h.serv.Create(r.Context(), service.CreateUserDTO(request))
		if err != nil {
			return err
		}

		return httplib.WriteJSON(w, http.StatusCreated, httplib.JSON{"user": user})
	}, h.errorHandler("handler.User.Create"))
}

func (h *User) Update() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		userNickname, ok := mux.Vars(r)["userNickname"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing nickname")
		}

		var request struct {
			Nickname    *string `json:"nickname"`
			Email       *string `json:"email"`
			AvatarPath  *string `json:"avatarPath"`
			Description *string `json:"description"`

			NewPassword *string `json:"newPassword"`
			OldPassword *string `json:"oldPassword"`
		}

		if err := httplib.DecodeJSON(r, &request); err != nil {
			return err
		}

		user, err := h.serv.UpdateByNickname(r.Context(), userNickname, service.UpdateUserDTO{
			Nickname:    request.Nickname,
			Email:       request.Email,
			AvatarPath:  request.AvatarPath,
			Description: request.Description,
			NewPassword: request.NewPassword,
			OldPassword: request.OldPassword,
		})
		if err != nil {
			return err
		}

		return httplib.WriteJSON(w, http.StatusOK, httplib.JSON{"user": user})
	}, h.errorHandler("handler.User.Update"))
}

func (h *User) Delete() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		userNickname, ok := mux.Vars(r)["userNickname"]
		if !ok {
			return httplib.NewAPIError(http.StatusBadRequest, "missing nickname")
		}

		if err := h.serv.DeleteByNickname(r.Context(), userNickname); err != nil {
			return err
		}

		return httplib.NoContent(w)
	}, h.errorHandler("handler.User.Delete"))
}

func (h *User) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)

		if errors.Is(err, model.ErrUserNotFound) {
			err = httplib.NewAPIError(http.StatusNotFound, model.ErrUserNotFound.Error())
		}
		if errors.Is(err, model.ErrUserExists) {
			err = httplib.NewAPIError(http.StatusConflict, model.ErrUserExists.Error())
		}

		httplib.DefaultErrorHandler(w, r, err)
	}
}
