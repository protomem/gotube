package entity

import (
	"time"

	"github.com/google/uuid"
)

type ID = uuid.UUID

type User struct {
	ID ID `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Nickname       string `json:"nickname"`
	HashedPassword string `json:"-"`

	Email    string `json:"email"`
	Verified bool   `json:"isVerified"`

	AvatarPath  string `json:"avatarPath"`
	Description string `json:"description"`
}
