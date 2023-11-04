package http

import (
	"net/http"

	"github.com/protomem/gotube/internal/service"
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
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handl *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handl *AuthHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handl *AuthHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
