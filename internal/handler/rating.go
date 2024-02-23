package handler

import (
	"net/http"

	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

type Rating struct {
	logger logging.Logger
	serv   service.Rating
}

func NewRating(logger logging.Logger, serv service.Rating) *Rating {
	return &Rating{
		logger: logger.With("handler", "rating"),
		serv:   serv,
	}
}

func (h *Rating) Count() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "unimplemented"})
	}, h.errorHandler("handler.Rating.Count"))
}

func (h *Rating) Like() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "unimplemented"})
	}, h.errorHandler("handler.Rating.Like"))
}

func (h *Rating) Dislike() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "unimplemented"})
	}, h.errorHandler("handler.Rating.Dislike"))
}

func (h *Rating) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)
		httplib.DefaultErrorHandler(w, r, err)
	}
}
