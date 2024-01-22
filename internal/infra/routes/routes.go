package routes

import (
	"fmt"
	"net/http"
)

func Setup() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GoTube API Server v5.0.0")
	})

	return mux
}
