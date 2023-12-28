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

	mux.Use(app.authenticate)

	mux.Get("/status", app.handleStatus)

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/auth", func(mux chi.Router) {
			mux.Post("/register", app.handleRegister)
			mux.Post("/login", app.handleLogin)

			mux.Group(func(mux chi.Router) {
				mux.Use(app.requireAuthentication)

				mux.Delete("/logout", app.handleLogout)
				mux.Get("/refresh", app.handleRefreshToken)
			})
		})

		mux.Route("/users", func(mux chi.Router) {
			mux.Get("/{userNickname}", app.handleGetUser)

			mux.Group(func(mux chi.Router) {
				mux.Use(app.requireAuthentication)

				mux.Put("/{userNickname}", app.handleUpdateUser)
				mux.Patch("/{userNickname}", app.handleUpdateUser)
				mux.Delete("/{userNickname}", app.handleDeleteUser)
			})
		})

		mux.Route("/videos", func(mux chi.Router) {
			mux.Get("/new", app.handleGetNewVideos)
			mux.Get("/popular", app.handleGetPopularVideos)
			mux.Get("/search", app.handleSearchVideo)
			mux.Get("/{videoId}", app.handleGetVideo)

			mux.Group(func(mux chi.Router) {
				mux.Use(app.requireAuthentication)

				mux.Post("/", app.handleCreateVideo)
				mux.Put("/{videoId}", app.handleUpdateVideo)
				mux.Patch("/{videoId}", app.handleUpdateVideo)
				mux.Delete("/{videoId}", app.handleDeleteVideo)
			})

			mux.Route("/{videoId}/ratings", func(mux chi.Router) {
				mux.Get("/", app.handleGetStatRatings)

				mux.Group(func(mux chi.Router) {
					mux.Use(app.requireAuthentication)

					mux.Post("/", app.handleLike)
					mux.Delete("/", app.handleUnlike)
				})
			})

			mux.Route("/{videoId}/comments", func(mux chi.Router) {
				mux.Get("/", app.handleGetComments)

				mux.Group(func(mux chi.Router) {
					mux.Use(app.requireAuthentication)

					mux.Post("/", app.handleCreateComment)
					mux.Put("/{commentId}", app.handleUpdateComment)
					mux.Patch("/{commentId}", app.handleUpdateComment)
					mux.Delete("/{commentId}", app.handleDeleteComment)
				})
			})
		})

		mux.Route("/profile/{userNickname}", func(mux chi.Router) {
			mux.Route("/videos", func(mux chi.Router) {
				mux.Get("/", app.handleGetUserVideos)
				mux.Get("/search", app.handleSearchUserVideo)

				mux.Group(func(mux chi.Router) {
					mux.Use(app.requireAuthentication)

					mux.Get("/subs", app.handleGetVideosBySubscriptions)
				})
			})

			mux.Route("/subs", func(mux chi.Router) {
				mux.Get("/", app.handleGetStatSubscriptions)

				mux.Group(func(mux chi.Router) {
					mux.Use(app.requireAuthentication)

					mux.Post("/", app.handleSubscribe)
					mux.Delete("/", app.handleUnsubscribe)
				})
			})
		})
	})

	return mux
}
