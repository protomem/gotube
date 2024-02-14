package handler

import (
	"errors"
	"net/http"

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

		return httplib.WriteJSON(w, http.StatusOK, httplib.JSON{"user": user})
	}, h.errorHandler("handler.User.Create"))
}

func (h *User) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.Error("failed to handle request", "operation", op, "err", err)

		if errors.Is(err, model.ErrUserNotFound) {
			err = httplib.NewAPIError(http.StatusNotFound, model.ErrUserNotFound.Error())
		}
		if errors.Is(err, model.ErrUserExists) {
			err = httplib.NewAPIError(http.StatusConflict, model.ErrUserExists.Error())
		}

		httplib.DefaultErrorHandler(w, r, err)
	}
}
