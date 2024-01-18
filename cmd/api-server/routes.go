package main

import (
	"net/http"

	gorilla "github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	mux := gorilla.NewRouter()

	mux.NotFoundHandler = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowed)

	mux.Use(app.requestID)
	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)

	mux.Use(app.authenticate)

	mux.HandleFunc("/status", app.handleStatus).Methods(http.MethodGet)

	{
		mux := mux.PathPrefix("/users").Subrouter()

		mux.HandleFunc("/{userNickname}", app.handleGetUser).Methods(http.MethodGet)

		{
			mux := mux.NewRoute().Subrouter()
			mux.Use(app.requireAuthentication)

			mux.HandleFunc("/{userNickname}", app.handleUpdateUser).Methods(http.MethodPut, http.MethodPatch)
			mux.HandleFunc("/{userNickname}", app.handleDeleteUser).Methods(http.MethodDelete)
		}
	}

	{
		mux := mux.PathPrefix("/auth").Subrouter()

		mux.HandleFunc("/register", app.handleRegister).Methods(http.MethodPost)
		mux.HandleFunc("/login", app.handleLogin).Methods(http.MethodPost)

		{
			mux := mux.NewRoute().Subrouter()
			mux.Use(app.requireAuthentication)

			mux.HandleFunc("/refresh", app.handleRefreshToken).Methods(http.MethodGet)
			mux.HandleFunc("/logout", app.handleLogout).Methods(http.MethodDelete)
		}
	}

	{
		mux := mux.PathPrefix("/videos").Subrouter()

		mux.NewRoute().HandlerFunc(app.handleFindVideos).Methods(http.MethodGet)
		mux.HandleFunc("/user/{userNickname}", app.handleFindUserVideos).Methods(http.MethodGet)
		mux.HandleFunc("/{videoID}", app.handleGetVideo).Methods(http.MethodGet)

		{
			mux := mux.NewRoute().Subrouter()
			mux.Use(app.requireAuthentication)

			mux.HandleFunc("/user/{userNickname}/private", app.handleFindPrivateUserVideos).Methods(http.MethodGet)

			mux.NewRoute().HandlerFunc(app.handleCreateVideo).Methods(http.MethodPost)
			mux.HandleFunc("/{videoID}", app.handleUpdateVideo).Methods(http.MethodPut, http.MethodPatch)
			mux.HandleFunc("/{videoID}", app.handleDeleteVideo).Methods(http.MethodDelete)
		}
	}

	{
		mux := mux.PathPrefix("/comments").Subrouter()

		mux.HandleFunc("/video/{videoID}", app.handleFindComments).Methods(http.MethodGet)

		{
			mux := mux.NewRoute().Subrouter()
			mux.Use(app.requireAuthentication)

			mux.HandleFunc("/video/{videoID}", app.handleCreateComment).Methods(http.MethodPost)
			mux.HandleFunc("/{commentID}", app.handleDeleteComment).Methods(http.MethodDelete)
		}
	}

	return mux
}
