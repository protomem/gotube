package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	mux.NotFoundHandler = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowed)

	mux.Use(app.requestID)
	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)

	mux.HandleFunc("/status", app.handleStatus).Methods(http.MethodGet)

	mux.HandleFunc("/users/{userNickname}", app.handleGetUser).Methods(http.MethodGet)

	mux.HandleFunc("/auth/register", app.handleRegister).Methods(http.MethodPost)

	return mux
}
