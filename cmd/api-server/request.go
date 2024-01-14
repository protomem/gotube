package main

import (
	"net/http"

	"github.com/gorilla/mux"
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
