package database

import (
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = NewModelError(ErrNotFound, "user")
	ErrUserAlreadyExists = NewModelError(ErrAlreadyExists, "user")
)

type User struct {
	ID uuid.UUID `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	Nickname string `db:"nickname" json:"nickname"`
	Password string `db:"password" json:"-"`

	Email string `db:"email" json:"email"`

	AvatarPath  string `db:"avatar_path" json:"avatarPath"`
	Description string `db:"description" json:"description"`
}
