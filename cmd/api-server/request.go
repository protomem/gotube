package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/protomem/gotube/internal/domain/model"
)

func getURLParamFromRequest(r *http.Request, name string) (string, bool) {
	vars := mux.Vars(r)
	if value, ok := vars[name]; ok {
		return value, true
	}
	return "", false
}

func mustGetURLParamFromRequest(r *http.Request, name string) string {
	value, _ := getURLParamFromRequest(r, name)
	return value
}

func mustGetUserNicknameFromRequest(r *http.Request) string {
	return mustGetURLParamFromRequest(r, "userNickname")
}

func getVideoIDFromRequest(r *http.Request) (model.ID, bool) {
	value, ok := getURLParamFromRequest(r, "videoID")
	if !ok {
		return model.ID{}, false
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return model.ID{}, false
	}
	return id, true
}

func getCommentIDFromRequest(r *http.Request) (model.ID, bool) {
	value, ok := getURLParamFromRequest(r, "commentID")
	if !ok {
		return model.ID{}, false
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return model.ID{}, false
	}
	return id, true
}

func getHeaderValueFromRequest(r *http.Request, name string) (string, bool) {
	if value := r.Header.Get(name); value != "" {
		return value, true
	}
	return "", false
}

func getQueryValueFromRequest(r *http.Request, name string) (string, bool) {
	if r.URL.Query().Has(name) {
		return r.URL.Query().Get(name), true
	}
	return "", false
}

func getRefreshTokenFromRequest(r *http.Request) (string, bool) {
	token, ok := getHeaderValueFromRequest(r, HeaderXRefreshToken)
	if !ok {
		token, ok = getQueryValueFromRequest(r, "refresh_token")
	}

	return token, ok
}
