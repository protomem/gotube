package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrVideoNotFound      = errors.New("video not found")
	ErrVideoAlreadyExists = errors.New("video already exists")

	ErrCommentNotFound = errors.New("comment not found")
)

type ID = uuid.UUID

type User struct {
	ID        ID        `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Nickname string `json:"nickname"`
	Password string `json:"-"`

	Email    string `json:"email"`
	Verified bool   `json:"isVerified"`

	AvatarPath  string `json:"avatarPath"`
	Description string `json:"description"`
}

type PairTokens struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}

type Video struct {
	ID        ID        `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Title       string `json:"title"`
	Description string `json:"description"`

	ThumbnailPath string `json:"thumbnailPath"`
	VideoPath     string `json:"videoPath"`

	Author User `json:"author"`

	Public bool   `json:"isPublic"`
	Views  uint64 `json:"views"`
}

type Comment struct {
	ID        ID        `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Message string `json:"message"`

	Author  User `json:"author"`
	VideoID ID   `json:"videoId"`
}
