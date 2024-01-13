package ctxstore

import (
	"context"
	"log/slog"
	"net/http"
)

type Key string

const (
	_requestIDKey = Key("requestId")
	_loggerKey    = Key("logger")
)

func WithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, _requestIDKey, rid)
}

func RequestIDWrapRequest(r *http.Request, rid string) *http.Request {
	return r.WithContext(WithRequestID(r.Context(), rid))
}

func RequestID(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(_requestIDKey).(string)
	return rid, ok
}

func MustRequestID(ctx context.Context) string {
	rid, _ := RequestID(ctx)
	return rid
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, _loggerKey, logger)
}

func LoggerWrapRequest(r *http.Request, logger *slog.Logger) *http.Request {
	return r.WithContext(WithLogger(r.Context(), logger))
}

func Logger(ctx context.Context) (*slog.Logger, bool) {
	logger, ok := ctx.Value(_loggerKey).(*slog.Logger)
	return logger, ok
}

func MustLogger(ctx context.Context) *slog.Logger {
	logger, _ := Logger(ctx)
	return logger
}
