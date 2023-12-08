package model

import "github.com/google/uuid"

type ID = uuid.UUID

type User struct {
	ID        ID     `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	Nickname string `json:"nickname"`
	Password string `json:"-"`

	Email    string `json:"email"`
	Verified bool   `json:"isVerified"`

	AvatarPath  string `json:"avatarPath"`
	Description string `json:"description"`
}
