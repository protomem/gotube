package ctxstore

import (
	"context"
	"net/http"

	"github.com/protomem/gotube/internal/model"
)

type Key string

const (
	_traceID = Key("traceId")
	_user    = Key("user")
)

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, _traceID, traceID)
}

func RequestWithTraceID(r *http.Request, traceID string) *http.Request {
	return r.WithContext(WithTraceID(r.Context(), traceID))
}

func TraceID(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(_traceID).(string)
	return traceID, ok
}

func MustTraceID(ctx context.Context) string {
	traceID, _ := TraceID(ctx)
	return traceID
}

func WithUser(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, _user, user)
}

func RequestWithUser(r *http.Request, user model.User) *http.Request {
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
