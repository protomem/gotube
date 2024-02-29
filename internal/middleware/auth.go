package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/internal/service"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
)

type Auth struct {
	logger logging.Logger
	serv   service.Auth
}

func NewAuth(logger logging.Logger, serv service.Auth) *Auth {
	return &Auth{
		logger: logger.With("middleware", "auth"),
		serv:   serv,
	}
}

func (m *Auth) Authenticate() mux.MiddlewareFunc {
	return mux.MiddlewareFunc(httplib.NewMiddlewareFunc(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader, ok := r.Header[httplib.HeaderAuthorization]
			if !ok || len(authHeader) == 0 {
				m.logger.Debug("authorization header is empty")
				next(w, r)
				return
			}

			headerParts := strings.Split(authHeader[0], " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				m.logger.Error("invalid authorization header")
				_ = httplib.WriteJSON(w, http.StatusBadRequest, httplib.JSON{"message": "invalid authorization header"})
				return
			}

			token := headerParts[1]
			user, err := m.serv.Verify(r.Context(), token)
			if err != nil {
				m.logger.Error("invalid token")
				_ = httplib.WriteJSON(w, http.StatusUnauthorized, httplib.JSON{"message": "invalid token"})
				return
			}

			wr := ctxstore.RequestWithUser(r, user)

			next(w, wr)
		}
	}))
}

func (m *Auth) Protect() mux.MiddlewareFunc {
	return mux.MiddlewareFunc(httplib.NewMiddlewareFunc(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, ok := ctxstore.User(r.Context())
			if !ok {
				_ = httplib.WriteJSON(w, http.StatusForbidden, httplib.JSON{"message": "access denied"})
				return
			}

			next(w, r)
		}
	}))
}
