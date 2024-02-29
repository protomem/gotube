package handler

import (
	"errors"
	"net/http"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

type Auth struct {
	logger logging.Logger
	serv   service.Auth
}

func NewAuth(logger logging.Logger, serv service.Auth) *Auth {
	return &Auth{
		logger: logger.With("handler", "auth"),
		serv:   serv,
	}
}

func (h *Auth) Login() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := httplib.DecodeJSON(r, &request); err != nil {
			return err
		}

		token, user, err := h.serv.Login(r.Context(), service.LoginDTO(request))
		if err != nil {
			return err
		}

		return httplib.WriteJSON(w, http.StatusOK, httplib.JSON{
			"accesssToken": token,
			"user":         user,
		})
	}, h.errorHandler("handler.Auth.Login"))
}

func (h *Auth) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)

		if errors.Is(err, model.ErrUserNotFound) {
			err = httplib.NewAPIError(http.StatusNotFound, model.ErrUserNotFound.Error())
		}

		httplib.DefaultErrorHandler(w, r, err)
	}
}
