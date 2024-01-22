package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/protomem/gotube/internal/infra/handler"
	"github.com/protomem/gotube/pkg/logging"
)

func Setup(
	logger *logging.Logger,
	common *handler.Common,
) http.Handler {
	mux := chi.NewRouter()

	mux.Use(common.TraceID)
	mux.Use(common.LogAccess(logger))
	mux.Use(common.RecoverPanic)

	mux.NotFound(common.NotFound)
	mux.MethodNotAllowed(common.MethodNotAllowed)

	mux.Get("/status", common.HandleStatus)

	return mux
}
