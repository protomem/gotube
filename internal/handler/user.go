package handler

import (
	"net/http"

	"github.com/protomem/gotube/internal/access"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
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
		return response.Send(w, http.StatusOK, response.JSON{"user": "some_user"})
	})
}

func (h *User) Create() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusCreated, response.JSON{"user": "some_user"})
	})
}

func (h *User) Update() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{"user": "some_user"})
	})
}

func (h *User) Delete() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusNoContent, nil)
	})
}

func (h *User) apiFunc(apiFn response.APIFunc) http.HandlerFunc {
	return response.BuildHandlerFunc(h.errorHandler())(apiFn)
}

func (h *User) errorHandler() response.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		// TODO: Implement this
	}
}
