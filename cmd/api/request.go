package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/protomem/gotube/internal/cookies"
)

func getRefreshTokenFromRequest(r *http.Request, cookieSecret string) string {
	var token string
	if r.Header.Get("X-Refresh-Token") != "" {
		token = r.Header.Get("X-Refresh-Token")
	} else if r.URL.Query().Has("refresh_token") {
		token = r.URL.Query().Get("refresh_token")
	} else {
		sessionToken, err := cookies.ReadSigned(r, "session", cookieSecret)
		if err != nil {
			return ""
		}
		token = sessionToken
	}
	return token
}

func getUserNicknameFromRequest(r *http.Request) string {
	return chi.URLParam(r, "userNickname")
}
