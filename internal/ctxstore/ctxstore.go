package ctxstore

import (
	"context"
	"net/http"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/pkg/logging"
)

type Key string

const (
	_loggerKey  = Key("logger")
	_traceIDKey = Key("traceID")
	_userKey    = Key("user")
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

func WithTraceID(ctx context.Context, tid string) context.Context {
	return context.WithValue(ctx, _traceIDKey, tid)
}

func RequestWithTraceID(r *http.Request, tid string) *http.Request {
	return r.WithContext(WithTraceID(r.Context(), tid))
}

func TraceID(ctx context.Context) (string, bool) {
	tid, ok := ctx.Value(_traceIDKey).(string)
	if !ok {
		return "", false
	}
	return tid, true
}

func MustTraceID(ctx context.Context) string {
	tid, _ := ctx.Value(_traceIDKey).(string)
	return tid
}

func WithUser(ctx context.Context, user entity.User) context.Context {
	return context.WithValue(ctx, _userKey, user)
}

func RequestWithUser(r *http.Request, user entity.User) *http.Request {
	return r.WithContext(WithUser(r.Context(), user))
}

func User(ctx context.Context) (entity.User, bool) {
	u, ok := ctx.Value(_userKey).(entity.User)
	if !ok {
		return entity.User{}, false
	}
	return u, true
}

func MustUser(ctx context.Context) entity.User {
	u, _ := ctx.Value(_userKey).(entity.User)
	return u
}
