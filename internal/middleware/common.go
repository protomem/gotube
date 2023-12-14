package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/protomem/gotube/pkg/header"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
)

type Common struct {
	logger logging.Logger
}

func NewCommon(logger logging.Logger) *Common {
	return &Common{logger: logger}
}

func (m *Common) CORS() MiddlewareFunc {
	opts := cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{
			http.MethodGet, http.MethodPost,
			http.MethodPut, http.MethodPatch,
			http.MethodDelete, http.MethodOptions,
		},
		AllowedHeaders: []string{
			header.Accept, header.Authorization,
			header.ContentType, header.XCSRFToken,
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}

	return cors.Handler(opts)
}

func (m *Common) RealIP() MiddlewareFunc {
	return middleware.RealIP
}

func (m *Common) RequestLogging() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := m.logger.With(requestid.Key, requestid.Extract(ctx))

			ww := newResponseWrapper(w)

			begin := time.Now()
			next.ServeHTTP(ww, r)
			end := time.Now()

			logger.Info(
				"incoming request",
				"method", r.Method,
				"path", r.URL.String(),
				"begin", begin,
				"end", end,
				"ip", r.RemoteAddr,
				"duration", end.Sub(begin).String(),
				"status", ww.StatusCode,
				"bytes", ww.WrittenBytes,
			)
		})
	}
}

// TODO: Implement it
func (m *Common) Recoverer() MiddlewareFunc {
	return middleware.Recoverer
}

func (m *Common) CleanPath() MiddlewareFunc {
	return middleware.CleanPath
}

func (m *Common) StripSlashes() MiddlewareFunc {
	return middleware.StripSlashes
}

type responseWrapper struct {
	w http.ResponseWriter

	StatusCode   int
	WrittenBytes int
}

func newResponseWrapper(w http.ResponseWriter) *responseWrapper {
	return &responseWrapper{
		w: w,
	}
}

func (rw *responseWrapper) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWrapper) Write(b []byte) (int, error) {
	if rw.StatusCode == 0 {
		rw.StatusCode = http.StatusOK
	}

	rw.WrittenBytes += len(b)

	return rw.w.Write(b)
}

func (rw *responseWrapper) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.w.WriteHeader(statusCode)
}
