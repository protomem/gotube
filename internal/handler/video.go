package handler

import (
	"net/http"

	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

type Video struct {
	logger logging.Logger
	serv   service.Video
}

func NewVideo(logger logging.Logger, serv service.Video) *Video {
	return &Video{
		logger: logger.With("handler", "video"),
		serv:   serv,
	}
}

func (h *Video) Get() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "not implemented"})
	}, h.errorHandler("handler.Video.Get"))
}

func (h *Video) Creaate() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "not implemented"})
	}, h.errorHandler("handler.Video.Create"))
}

func (h *Video) Update() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "not implemented"})
	}, h.errorHandler("handler.Video.Update"))
}

func (h *Video) Delete() http.HandlerFunc {
	return httplib.NewEndpointWithErroHandler(func(w http.ResponseWriter, r *http.Request) error {
		return httplib.WriteJSON(w, http.StatusInternalServerError, httplib.JSON{"message": "not implemented"})
	}, h.errorHandler("handler.Video.Delete"))
}

func (h *Video) errorHandler(op string) httplib.ErroHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.WithContext(r.Context()).Error("failed to handle request", "operation", op, "err", err)
		httplib.DefaultErrorHandler(w, r, err)
	}
}
