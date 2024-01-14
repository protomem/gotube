package model

import (
	"time"

	"github.com/google/uuid"
)

type ID = uuid.UUID

var (
	ErrUserNotFound      = NewError("user", ErrNotFound)
	ErrUserAlreadyExists = NewError("user", ErrAlreadyExists)
)

type User struct {
	ID uuid.UUID `json:"id" db:"id"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`

	Nickname string `json:"nickname" db:"nickname"`
	Password string `json:"-" db:"password"`

	Email string `json:"email" db:"email"`

	AvatarPath  string `json:"avatarPath" db:"avatar_path"`
	Description string `json:"description" db:"description"`
}

var ErrSessionNotFound = NewError("session", ErrNotFound)

type Session struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
	UserID ID        `json:"userId"`
}
