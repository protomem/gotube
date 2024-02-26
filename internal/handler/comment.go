package handler

import (
	"net/http"

	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

type Comment struct {
	logger logging.Logger
	serv   service.Comment
}

func NewComment(logger logging.Logger, serv service.Comment) *Comment {
	return &Comment{
		logger: logger.With("handler", "comment"),
		serv:   serv,
	}
}

func (h *Comment) List() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "unimplemented"})
	}, h.errorHandler("handler.Comment.List"))
}

func (h *Comment) Create() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "unimplemented"})
	}, h.errorHandler("handler.Comment.List"))
}

func (h *Comment) Delete() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "unimplemented"})
	}, h.errorHandler("handler.Comment.List"))
}

func (h *Comment) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)
		httplib.DefaultErrorHandler(w, r, err)
	}
}
