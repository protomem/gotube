package database

import (
	"time"

	"github.com/google/uuid"
)

type UserTable struct {
	ID uuid.UUID `db:"id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Nickname       string `db:"nickname"`
	HashedPassword string `db:"hashed_password"`

	Email    string `db:"email"`
	Verified bool   `db:"is_verified"`

	AvatarPath  string `db:"avatar_path"`
	Description string `db:"description"`
}

type UserDAO struct {
	db *DB
}

func (db *DB) UserDAO() *UserDAO {
	return &UserDAO{db}
}
