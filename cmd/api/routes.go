package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Use(app.CORS)

	mux.Use(app.cleanPath)
	mux.Use(app.stripSlashes)

	mux.Use(app.requestID)
	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)

	mux.Get("/status", app.handleStatus)

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/auth", func(mux chi.Router) {
			mux.Post("/register", app.handleRegister)
			mux.Post("/login", app.handleLogin)

			mux.Group(func(mux chi.Router) {
				mux.Delete("/logout", app.handleLogout)
				mux.Get("/refresh", app.handleRefreshToken)
			})
		})

		mux.Route("/users", func(mux chi.Router) {
			mux.Get("/{userNickname}", app.handleGetUser)

			mux.Group(func(mux chi.Router) {
				mux.Patch("/{userNickname}", app.handleUpdateUser)
				mux.Delete("/{userNickname}", app.handleDeleteUser)
			})
		})
	})

	return mux
}
