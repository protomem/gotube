package model

import (
	"time"

	"github.com/google/uuid"
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
