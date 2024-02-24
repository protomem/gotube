package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ID = uuid.UUID

type Model struct {
	ID        ID        `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type User struct {
	Model

	Nickname string `json:"nickname"`
	Password string `json:"-"`

	Email    string `json:"email"`
	Verified bool   `json:"isVerified"`

	AvatarPath  string `json:"avatarPath"`
	Description string `json:"description"`
}

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrSubscriptionExists   = errors.New("subscription already exists")
)

type Subscription struct {
	Model

	FromUserID ID `json:"fromUser"`
	ToUserID   ID `json:"toUser"`
}

var (
	ErrVideoNotFound = errors.New("video not found")
	ErrVideoExists   = errors.New("video already exists")
)

type Video struct {
	Model

	Title       string `json:"title"`
	Description string `json:"description"`

	ThumbnailPath string `json:"thumbnailPath"`
	VideoPath     string `json:"videoPath"`

	Author User `json:"author"`

	Public bool  `json:"isPublic"`
	Views  int64 `json:"views"`
}

type Rating struct {
	Model

	UserID  ID `json:"userId"`
	VideoID ID `json:"videoId"`

	Like bool `json:"isLike"`
}
