package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Setup() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GoTube API Server v5.0")
	})

	return mux
}
