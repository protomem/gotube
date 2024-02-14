package ctxstore

import (
	"context"
	"net/http"
)

type Key string

const (
	_traceID = Key("traceId")
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
