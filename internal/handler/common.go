package handler

import (
	"net/http"

	"github.com/protomem/gotube/pkg/httplib"
)

type Common struct{}

func NewCommon() *Common {
	return &Common{}
}

func (*Common) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (*Common) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (*Common) Health() http.HandlerFunc {
	return httplib.NewEndpoint(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusOK, httplib.JSON{"status": "OK"})
	})
}
