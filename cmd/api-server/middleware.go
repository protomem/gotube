package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/internal/domain/jwt"
	"github.com/protomem/gotube/internal/domain/usecase"
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
		var (
			ip     = realip.FromRequest(r)
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
			rid    = ctxstore.MustRequestID(r.Context())
		)

		mw := response.NewMetricsResponseWriter(w)
		wr := ctxstore.LoggerWrapRequest(r, app.logger.With(slog.String("traceId", rid)))

		next.ServeHTTP(mw, wr)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto, "traceId", rid)
		responseAttrs := slog.Group("repsonse", "status", mw.StatusCode, "size", mw.BytesCount)

		app.logger.Info("access", userAttrs, requestAttrs, responseAttrs)
	})
}

func (app *application) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rid string
		if r.Header.Get(HeaderXRequestID) != "" {
			rid = r.Header.Get(HeaderXRequestID)
		} else {
			id, err := uuid.NewRandom()
			if err != nil {
				id = uuid.Nil
			}
			rid = id.String()
		}

		wr := ctxstore.RequestIDWrapRequest(r, rid)
		next.ServeHTTP(w, wr)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := getHeaderValueFromRequest(r, HeaderAuthorization)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 && headerParts[0] != "Bearer" {
			app.errorMessage(w, r, http.StatusBadRequest, "invalid authorization header", nil)
			return
		}

		authToken := headerParts[1]
		user, err := usecase.VerifyToken(app.config.auth.secret, app.db).Invoke(r.Context(), authToken)
		if err != nil {
			if errors.Is(err, jwt.ErrInvalidToken) {
				app.errorMessage(w, r, http.StatusUnauthorized, jwt.ErrInvalidToken.Error(), nil)
				return
			}

			app.serverError(w, r, err)
			return
		}

		wr := ctxstore.UserWrapRequest(r, user)
		next.ServeHTTP(w, wr)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := ctxstore.User(r.Context())
		if !ok {
			app.errorMessage(w, r, http.StatusForbidden, "access denied", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
