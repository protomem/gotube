package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/protomem/gotube/internal/cookies"
)

func getAccessTokenFromRequest(r *http.Request) string {
	header := r.Header.Get(HeaderAuthorization)
	hedaerParts := strings.Split(header, " ")
	if len(hedaerParts) != 2 && hedaerParts[0] != "Bearer" {
		return ""
	}
	return hedaerParts[1]
}

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
