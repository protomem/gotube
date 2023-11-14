package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")

	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrSubscriptionExists   = errors.New("subscription already exists")

	ErrVideoNotFound = errors.New("video not found")
	ErrVideoExists   = errors.New("video already exists")

	ErrRatingNotFound = errors.New("rating not found")
	ErrRatingExists   = errors.New("rating already exists")
)

type User struct {
	ID uuid.UUID `json:"id"`

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

type Subscription struct {
	ID uuid.UUID `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	FromUser User `json:"fromUser"`
	ToUser   User `json:"toUser"`
}

type Video struct {
	ID uuid.UUID `json:"id"`

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

type Rating struct {
	ID uuid.UUID `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Like bool `json:"isLike"`

	VideoID uuid.UUID `json:"videoID"`
	UserID  uuid.UUID `json:"userID"`
}

type Comment struct {
	ID uuid.UUID `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Content string `json:"content"`

	Author User `json:"author"`

	VideoID uuid.UUID `json:"videoID"`
}
