package database

import (
	"context"
	"strconv"
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

func (db *DB) GetUserByNickname(ctx context.Context, nickname string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT * FROM users WHERE nickname = $1 LIMIT 1`
	args := []any{nickname}

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

func (db *DB) GetUserByEmail(ctx context.Context, email string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`
	args := []any{email}

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

type UpdateUserDTO struct {
	Nickname    *string
	Password    *string
	Email       *string
	AvatarPath  *string
	Description *string
}

func (db *DB) UpdateUser(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	counter := 1
	query := `UPDATE users SET updated_at = now()`
	args := []any{id}

	if dto.Nickname != nil {
		counter++
		query += `, nickname = $` + strconv.Itoa(counter)
		args = append(args, *dto.Nickname)
	}
	if dto.Password != nil {
		counter++
		query += `, password = $` + strconv.Itoa(counter)
		args = append(args, *dto.Password)
	}
	if dto.Email != nil {
		counter++
		query += `, email = $` + strconv.Itoa(counter)

	}
	if dto.AvatarPath != nil {
		counter++
		query += `, avatar_path = $` + strconv.Itoa(counter)
		args = append(args, *dto.AvatarPath)
	}
	if dto.Description != nil {
		counter++
		query += `, description = $` + strconv.Itoa(counter)
		args = append(args, *dto.Description)
	}

	query += ` WHERE id = $1`

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
