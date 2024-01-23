package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserEntry struct {
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

func (dao *UserDAO) GetByID(ctx context.Context, id uuid.UUID) (UserEntry, error) {
	const op = "database.UserDAO.GetByID"

	query := `SELECT * FROM users WHERE id = $1 LIMIT 1`
	args := []any{id}

	var user UserEntry

	if err := dao.db.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		return UserEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (dao *UserDAO) GetByNickname(ctx context.Context, nickname string) (UserEntry, error) {
	const op = "database.UserDAO.GetByNickname"

	query := `SELECT * FROM users WHERE nickname = $1 LIMIT 1`
	args := []any{nickname}

	var user UserEntry

	if err := dao.db.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		return UserEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (dao *UserDAO) GetByEmail(ctx context.Context, email string) (UserEntry, error) {
	const op = "database.UserDAO.GetByEmail"

	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`
	args := []any{email}

	var user UserEntry

	if err := dao.db.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		return UserEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

type InsertUserDTO struct {
	Nickname string
	Password string
	Email    string
}

func (dao *UserDAO) Insert(ctx context.Context, dto InsertUserDTO) (uuid.UUID, error) {
	const op = "database.UserDAO.Insert"

	query := `INSERT INTO users(nickname, hashed_password, email) VALUES ($1, $2, $3) RETURNING id`
	args := []any{dto.Nickname, dto.Password, dto.Email}

	var id uuid.UUID

	if err := dao.db.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
