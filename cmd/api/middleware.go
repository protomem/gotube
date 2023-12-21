package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/response"

	"github.com/tomasen/realip"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) logAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := response.NewMetricsResponseWriter(w)
		next.ServeHTTP(mw, r)

		var (
			ip     = realip.FromRequest(r)
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request",
			"method", method,
			"url", url,
			"proto", proto,
			string(_contextRequestIDKey), contextGetRequestID(r),
		)
		responseAttrs := slog.Group("response", "status", mw.StatusCode, "size", mw.BytesCount)

		app.logger.Info("access", userAttrs, requestAttrs, responseAttrs)
	})
}

func (app *application) cleanPath(next http.Handler) http.Handler {
	return middleware.CleanPath(next)
}

func (app *application) stripSlashes(next http.Handler) http.Handler {
	return middleware.StripSlashes(next)
}

func (app *application) CORS(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{
			http.MethodGet, http.MethodPost,
			http.MethodPut, http.MethodPatch,
			http.MethodDelete, http.MethodOptions,
		},
		AllowedHeaders: []string{
			HeaderAccept, HeaderAuthorization,
			HeaderContentType, HeaderXCSRFToken,
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})(next)
}

func (app *application) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := uuid.New().String()

		r = contextSetRequestID(r, rid)
		w.Header().Set(HeaderXRequestID, rid)

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getAccessTokenFromRequest(r)
		if token == "" {
			app.logger.Debug("auth header missing or invalid")
			next.ServeHTTP(w, r)
			return
		}

		userID, err := jwt.Parse(token, jwt.ParseParams{
			SigningKey: app.config.auth.secretKey,
			Issuer:     app.config.baseURL,
		})
		if err != nil {
			app.errorMessage(w, r, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		user, err := app.db.GetUser(r.Context(), userID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				app.errorMessage(w, r, http.StatusUnauthorized, database.ErrUserNotFound.Error(), nil)
				return
			}

			app.serverError(w, r, err)
			return
		}

		r = contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := contextGetUser(r); !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
