package handler

import (
	"net/http"

	"github.com/protomem/gotube/internal/access"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/response"
)

type Video struct {
	logger logging.Logger
	serv   service.Video
	accmng access.Manager
}

func NewVideo(logger logging.Logger, serv service.Video, accmng access.Manager) *Video {
	return &Video{
		logger: logger.With("handler", "video", "handlerType", "http"),
		serv:   serv,
		accmng: accmng,
	}
}

func (h *Video) List() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{"videos": "some_videos"})
	})
}

func (h *Video) Get() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{"video": "some_video"})
	})
}

func (h *Video) Create() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusCreated, response.JSON{"video": "some_video"})
	})
}

func (h *Video) Update() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{"video": "some_video"})
	})
}

func (h *Video) Delete() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusNoContent, nil)
	})
}

func (h *Video) apiFunc(apiFn response.APIFunc) http.HandlerFunc {
	return response.BuildHandlerFunc(h.errorHandler())(apiFn)
}

func (h *Video) errorHandler() response.ErrorHandler {
	return response.DefaultErrorHandler(h.logger, "handler:User.errorHandler")
}
