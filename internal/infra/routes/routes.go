package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/response"
)

func Setup(logger *logging.Logger) http.Handler {
	mux := chi.NewRouter()

	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wr := ctxstore.RequestWithLogger(r, logger)
			next.ServeHTTP(w, wr)
		})
	})

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctxstore.MustLogger(r.Context()).Debug("incomig request")
		_ = response.JSON(w, http.StatusOK, response.Data{"status": "OK"})
	})

	return mux
}
