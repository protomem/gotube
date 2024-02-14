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

var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")

type User struct {
	Model

	Nickname string `json:"nickname"`
	Password string `json:"-"`

	Email    string `json:"email"`
	Verified bool   `json:"isVerified"`

	AvatarPath  string `json:"avatarPath"`
	Description string `json:"description"`
}
