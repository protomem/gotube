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
	user *handler.User,
	auth *handler.Auth,
) http.Handler {
	mux := chi.NewRouter()

	mux.Use(common.TraceID)
	mux.Use(common.LogAccess(logger))
	mux.Use(common.RecoverPanic)

	mux.NotFound(common.NotFound)
	mux.MethodNotAllowed(common.MethodNotAllowed)

	mux.Get("/status", common.HandleStatus)

	mux.Route("/users", func(mux chi.Router) {
		mux.Get("/{nickname}", user.HandleGet)
	})

	mux.Route("/auth", func(mux chi.Router) {
		mux.Post("/register", auth.HandleRegister)
		mux.Post("/login", auth.HandleLogin)
		mux.Get("/refresh", auth.HandleRefreshTokens)
		mux.Delete("/logout", auth.HandleLogout)
	})

	return mux
}
