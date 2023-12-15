package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	_contextRequestIDKey = contextKey("requestId")
)

func contextSetRequestID(r *http.Request, rid string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), _contextRequestIDKey, rid))
}

func contextGetRequestID(r *http.Request) string {
	rid, ok := r.Context().Value(_contextRequestIDKey).(string)
	if !ok || rid == "" {
		return uuid.Nil.String()
	}
	return rid
}
