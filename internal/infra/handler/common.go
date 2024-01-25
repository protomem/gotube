package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/response"
	"github.com/tomasen/realip"
)

type Common struct {
	*Base
}

func NewCommon() *Common {
	return &Common{NewBase()}
}

func (h *Common) HandleStatus(w http.ResponseWriter, r *http.Request) {
	h.MustSendJSON(w, r, http.StatusOK, response.Data{"status": "ok"})
}

func (h *Common) TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tid, _ := uuid.NewRandom()
		wr := ctxstore.RequestWithTraceID(r, tid.String())
		next.ServeHTTP(w, wr)
	})
}

func (h *Common) LogAccess(logger *logging.Logger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				ip     = realip.FromRequest(r)
				method = r.Method
				url    = r.URL.String()
				proto  = r.Proto
				tid    = ctxstore.MustTraceID(r.Context())
			)

			mw := response.NewMetricsResponseWriter(w)
			wr := ctxstore.RequestWithLogger(r, logger.With(slog.String("traceId", tid)))

			next.ServeHTTP(mw, wr)

			userAttrs := slog.Group("user", "ip", ip)
			requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto, "traceId", tid)
			responseAttrs := slog.Group("repsonse", "status", mw.StatusCode, "size", mw.BytesCount)

			logger.Info("access", userAttrs, requestAttrs, responseAttrs)
		})
	}
}

func (h *Common) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				h.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (h *Common) CleanPath(next http.Handler) http.Handler {
	return middleware.CleanPath(next)
}

func (h *Common) StripSlashes(next http.Handler) http.Handler {
	return middleware.StripSlashes(next)
}
