package ctxstore

import (
	"context"
	"net/http"

	"github.com/protomem/gotube/pkg/logging"
)

type Key string

const (
	_loggerKey = Key("logger")
)

func WithLogger(ctx context.Context, l *logging.Logger) context.Context {
	return context.WithValue(ctx, _loggerKey, l)
}

func RequestWithLogger(r *http.Request, l *logging.Logger) *http.Request {
	return r.WithContext(WithLogger(r.Context(), l))
}

func Logger(ctx context.Context) (*logging.Logger, bool) {
	l, ok := ctx.Value(_loggerKey).(*logging.Logger)
	if !ok {
		return nil, false
	}
	return l, true
}

func MustLogger(ctx context.Context) *logging.Logger {
	l, _ := ctx.Value(_loggerKey).(*logging.Logger)
	return l
}
