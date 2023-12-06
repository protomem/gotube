package handler

import (
	"encoding/json"
	"net/http"

	"github.com/protomem/gotube/pkg/header"
	"github.com/protomem/gotube/pkg/logging"
)

type Common struct {
	logger logging.Logger
}

func NewCommon(logger logging.Logger) *Common {
	return &Common{logger: logger}
}

func (h *Common) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(header.ContentType, header.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JSON{"message": "pong"})
	}
}
