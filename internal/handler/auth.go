package handler

import (
	"net/http"

	"github.com/protomem/gotube/internal/access"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
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
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusCreated, response.JSON{"user": "some_user"})
	})
}

func (h *Auth) Login() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{"user": "some_user"})
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
