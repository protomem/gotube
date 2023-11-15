package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httpheader"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type AuthMiddleware struct {
	logger logging.Logger
	serv   service.Auth
}

func NewAuthMiddleware(logger logging.Logger, serv service.Auth) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger.With("middleware", "auth", "middlewareType", "http"),
		serv:   serv,
	}
}

func (mdw *AuthMiddleware) Authenticate() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "http.AuthMiddleware.Authenticate"
			var err error

			ctx := r.Context()
			logger := mdw.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)

			defer func() {
				if err != nil {
					logger.Error("failed to send response", "error", err)
				}
			}()

			header := r.Header.Get(httpheader.Authorization)
			if header == "" {
				logger.Debug("no authorization header")

				next.ServeHTTP(w, r)

				return
			}

			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				logger.Error("invalid authorization header")

				w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
				w.WriteHeader(http.StatusUnauthorized)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid authorization header",
				})

				return
			}

			token := headerParts[1]
			_, payload, err := mdw.serv.VerifyToken(ctx, token)
			if err != nil {
				logger.Error("failed to verify token", "error", err)

				w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
				w.WriteHeader(http.StatusUnauthorized)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "failed to verify token",
				})

				return
			}

			ctx = jwt.Inject(ctx, payload)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func (mdw *AuthMiddleware) Protect() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "http.AuthMiddleware.Protect"
			var err error

			ctx := r.Context()
			logger := mdw.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)

			defer func() {
				if err != nil {
					logger.Error("failed to send response", "error", err)
				}
			}()

			_, ok := jwt.Extract(ctx)
			if !ok {
				logger.Error("auth payload missing")

				w.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)
				w.WriteHeader(http.StatusForbidden)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "not authenticated",
				})

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
