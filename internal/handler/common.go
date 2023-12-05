package handler

import (
	"encoding/json"
	"net/http"

	"github.com/protomem/gotube/pkg/logging"
)

type Common struct {
	logger logging.Logger
}

func NewCommon(logger logging.Logger) *Common {
	return &Common{logger: logger.With("handler", "common")}
}

func (h *Common) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(HeaderContentType, ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JSON{"message": "pong"})
	}
}
