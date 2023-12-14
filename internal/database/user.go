package database

import (
	"context"
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

func (db *DB) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT * FROM users WHERE id = $1 LIMIT 1`
	args := []any{id}

	var user User

	if err := db.
		QueryRowxContext(ctx, query, args...).
		StructScan(&user); err != nil {
		if IsNoRows(err) {
			return User{}, ErrUserNotFound
		}

		return User{}, err
	}

	return user, nil
}

type InsertUserDTO struct {
	Nickname string
	Email    string
	Password string
}

func (db *DB) InsertUser(ctx context.Context, dto InsertUserDTO) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `INSERT INTO users (nickname, email, password) VALUES ($1, $2, $3) RETURNING id`
	args := []any{dto.Nickname, dto.Email, dto.Password}

	var id uuid.UUID

	if err := db.
		QueryRowxContext(ctx, query, args...).
		Scan(&id); err != nil {
		if IsKeyConflict(err) {
			return uuid.Nil, ErrUserAlreadyExists
		}

		return uuid.Nil, err
	}

	return id, nil
}
