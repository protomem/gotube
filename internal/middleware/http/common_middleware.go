package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/requestid"
	"github.com/rs/cors"
)

type CommonMiddleware struct {
	logger logging.Logger
}

func NewCommonMiddleware(logger logging.Logger) *CommonMiddleware {
	return &CommonMiddleware{
		logger: logger.With("middleware", "common", "middlewareType", "http"),
	}
}

func (mdw *CommonMiddleware) RequestID() mux.MiddlewareFunc {
	return requestid.Middleware()
}

func (mdw *CommonMiddleware) Logging() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "http.CommonMiddleware.Logging"

			ctx := r.Context()
			logger := mdw.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)

			ww := newResponseWrapper(w)

			begin := time.Now()
			next.ServeHTTP(ww, r)
			end := time.Now()

			logger.Info(
				"incoming request",
				"method", r.Method,
				"url", r.URL.String(),
				"begin", begin,
				"end", end,
				"duration", end.Sub(begin).String(),
				"status", ww.StatusCode,
				"bytes", ww.WrittenBytes,
			)
		})
	}
}

func (mdw *CommonMiddleware) Recovery() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "http.CommonMiddleware.Recovery"

			ctx := r.Context()
			logger := mdw.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)

			defer func() {
				err := recover()
				if err != nil {
					logger.Error("panic recovered", "error", err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func (mdw *CommonMiddleware) CORS() mux.MiddlewareFunc {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodOptions, http.MethodGet, http.MethodPost,
			http.MethodPatch, http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler
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
