package handler

import (
	"net/http"

	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

type Media struct {
	logger logging.Logger
}

func NewMedia(logger logging.Logger) *Media {
	return &Media{
		logger: logger.With("handler", "media"),
	}
}

func (h *Media) Get() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "not implemented"})
	}, h.errorHandler("handler.Media.Get"))
}

func (h *Media) Create() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "not implemented"})
	}, h.errorHandler("handler.Media.Get"))
}

func (h *Media) Delete() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "not implemented"})
	}, h.errorHandler("handler.Media.Get"))
}

func (h *Media) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)
		httplib.DefaultErrorHandler(w, r, err)
	}
}
