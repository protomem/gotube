package http

import (
	"encoding/json"
	"net/http"

	"github.com/protomem/gotube/pkg/httpheader"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type CommonHandler struct {
	logger logging.Logger
}

func NewCommonHandler(logger logging.Logger) *CommonHandler {
	return &CommonHandler{
		logger: logger.With("handler", "common", "handlerType", "http"),
	}
}

func (handl *CommonHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.CommonHandler.Ping"

		ctx := r.Context()
		logger := handl.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(map[string]string{
			"message": "pong",
		})
		if err != nil {
			logger.Error("failed to send response", "error", err)

			return
		}
	}
}
