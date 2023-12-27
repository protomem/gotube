package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/cookies"
	"github.com/protomem/gotube/internal/database"
)

const (
	_defaultLimit  = 10
	_defaultOffset = 0
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

func getVideoIDFromRequest(r *http.Request) (uuid.UUID, error) {
	videoIDRaw := chi.URLParam(r, "videoId")
	videoID, err := uuid.Parse(videoIDRaw)
	if err != nil {
		return uuid.Nil, err
	}
	return videoID, nil
}

func getFindOptionsFromRequest(r *http.Request) (database.FindOptions, error) {
	opts := database.FindOptions{
		Limit:  _defaultLimit,
		Offset: _defaultOffset,
	}
	if r.URL.Query().Has("limit") {
		limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			return database.FindOptions{}, err
		}
		opts.Limit = limit
	}
	if r.URL.Query().Has("offset") {
		offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
		if err != nil {
			return database.FindOptions{}, err
		}
		opts.Offset = offset
	}
	return opts, nil
}

func getSearchQueryFromRequest(r *http.Request) string {
	return r.URL.Query().Get("q")
}

func getCommentIDFromRequest(r *http.Request) (uuid.UUID, error) {
	commentIDRaw := chi.URLParam(r, "commentId")
	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return uuid.Nil, err
	}
	return commentID, nil
}
