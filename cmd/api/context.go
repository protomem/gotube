package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/database"
)

type contextKey string

const (
	_contextRequestIDKey = contextKey("requestId")
	_contextUserKey      = contextKey("user")
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

func contextSetUser(r *http.Request, user database.User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), _contextUserKey, user))
}

func contextGetUser(r *http.Request) (database.User, bool) {
	user, ok := r.Context().Value(_contextUserKey).(database.User)
	if !ok {
		return database.User{}, false
	}
	return user, true
}
