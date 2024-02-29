package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/ctxstore"
	"github.com/protomem/gotube/pkg/httplib"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/tomasen/realip"
)

type Common struct{}

func NewCommon() *Common {
	return &Common{}
}

func (m *Common) TraceID() mux.MiddlewareFunc {
	return mux.MiddlewareFunc(httplib.NewMiddlewareFunc(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			traceID := m.generateTraceID()
			wr := ctxstore.RequestWithTraceID(r, traceID)
			next(w, wr)
		}
	}))
}

func (*Common) LogAccess(logger logging.Logger) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(httplib.NewMiddlewareFunc(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var (
				ip     = realip.FromRequest(r)
				method = r.Method
				url    = r.URL.String()
				proto  = r.Proto
			)

			mw := httplib.NewMetricsResponseWriter(w)

			next(mw, r)

			logger.WithContext(r.Context()).Info(
				"incoming request",
				"method", method,
				"url", url,
				"statusCode", mw.StatusCode,
				"size", mw.BytesCount,
				"proto", proto,
				"ip", ip,
			)
		}
	}))
}

func (*Common) Recovery(logger logging.Logger) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(httplib.NewMiddlewareFunc(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered", "error", err)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next(w, r)
		}
	}))
}

func (*Common) generateTraceID() string {
	id, _ := uuid.NewRandom()
	return id.String()
}
