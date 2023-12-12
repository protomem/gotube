package handler

import (
	"net/http"

	"github.com/protomem/gotube/internal/access"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/response"
)

type Comment struct {
	logger logging.Logger
	serv   service.Comment
	accmng access.Manager
}

func NewComment(logger logging.Logger, serv service.Comment, accmng access.Manager) *Comment {
	return &Comment{
		logger: logger.With("handler", "comment", "handlerType", "http"),
		serv:   serv,
		accmng: accmng,
	}
}

func (h *Comment) List() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{
			"videos": "some_videos",
		})
	})
}

func (h *Comment) Create() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{
			"comment": "some_comment",
		})
	})
}

func (h *Comment) Update() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{
			"comment": "some_comment",
		})
	})
}

func (h *Comment) Delete() http.HandlerFunc {
	return h.apiFunc(func(w http.ResponseWriter, r *http.Request) error {
		return response.Send(w, http.StatusOK, response.JSON{
			"comment": "some_comment",
		})
	})
}

func (h *Comment) apiFunc(apiFn response.APIFunc) http.HandlerFunc {
	return response.BuildHandlerFunc(h.errorHandler())(apiFn)
}

func (h *Comment) errorHandler() response.ErrorHandler {
	return response.DefaultErrorHandler(h.logger, "handler:Auth.errorHandler")
}
