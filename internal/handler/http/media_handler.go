package http

import (
	"net/http"

	"github.com/protomem/gotube/internal/storage"
	"github.com/protomem/gotube/pkg/logging"
)

type MediaHandler struct {
	logger logging.Logger
	store  storage.Storage
}

func NewMediaHandler(logger logging.Logger, store storage.Storage) *MediaHandler {
	return &MediaHandler{
		logger: logger.With("handler", "media", "handlerType", "http"),
		store:  store,
	}
}

func (handl *MediaHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handl *MediaHandler) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handl *MediaHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
