package middleware

import (
	"net/http"
	"strings"

	"github.com/protomem/gotube/internal/jwt"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/header"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
	"github.com/protomem/gotube/pkg/response"
)

type Auth struct {
	logger logging.Logger
	serv   service.Auth
}

func NewAuth(logger logging.Logger, serv service.Auth) *Auth {
	return &Auth{
		logger: logger.With("middleware", "auth", "middlewareType", "http"),
		serv:   serv,
	}
}

func (m *Auth) Authenticate() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware:Auth.Authenticate"

			ctx := r.Context()
			logger := m.logger.With(
				"operation", op,
				requestid.Key, requestid.Extract(ctx),
			)

			authHeader := r.Header.Get(header.Authorization)
			if authHeader == "" {
				logger.Debug("authorization header is missing")

				next.ServeHTTP(w, r)

				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 && headerParts[0] != "Bearer" {
				logger.Error("authorization header is invalid")

				_ = response.Send(w, http.StatusUnauthorized, response.JSON{
					"error": "authorization header is invalid",
				})

				return
			}

			token := headerParts[1]
			_, authPayload, err := m.serv.VerifyToken(ctx, token)
			if err != nil {
				logger.Error("failed to verify token", "error", err)

				_ = response.Send(w, http.StatusUnauthorized, response.JSON{
					"error": "failed to verify token",
				})

				return
			}

			ctxWithAuthPayload := jwt.Inject(ctx, authPayload)
			next.ServeHTTP(w, r.WithContext(ctxWithAuthPayload))
		})
	}
}

func (m *Auth) IsAuthenticated() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware:Auth.IsAuthenticated"

			ctx := r.Context()
			logger := m.logger.With(
				"operation", op,
				requestid.Key, requestid.Extract(ctx),
			)

			_, ok := jwt.Extract(ctx)
			if !ok {
				logger.Error("auth payload missing")

				_ = response.Send(w, http.StatusForbidden, nil)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
