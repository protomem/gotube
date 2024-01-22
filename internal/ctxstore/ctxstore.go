package ctxstore

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/protomem/gotube/internal/domain/model"
)

type Key string

const (
	_requestIDKey = Key("requestId")
	_loggerKey    = Key("logger")
	_user         = Key("user")
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

func WithUser(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, _user, user)
}

func UserWrapRequest(r *http.Request, user model.User) *http.Request {
	return r.WithContext(WithUser(r.Context(), user))
}

func User(ctx context.Context) (model.User, bool) {
	user, ok := ctx.Value(_user).(model.User)
	return user, ok
}

func MustUser(ctx context.Context) model.User {
	user, _ := User(ctx)
	return user
}
