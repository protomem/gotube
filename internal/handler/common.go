package handler

import (
	"net/http"

	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/response"
)

type Common struct {
	logger logging.Logger
}

func NewCommon(logger logging.Logger) *Common {
	return &Common{logger: logger}
}

func (h *Common) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_ = response.Send(w, http.StatusOK, response.JSON{"message": "pong"})
	}
}
