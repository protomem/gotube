package http

import (
	"net/http"

	"github.com/protomem/gotube/internal/service"
)

type UserHandler struct {
	serv service.User
}

func NewUserHandler(serv service.User) *UserHandler {
	return &UserHandler{
		serv: serv,
	}
}

func (handl *UserHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handl *UserHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handl *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handl *UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
