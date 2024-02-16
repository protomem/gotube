package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/database/sqlite"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.User = (*User)(nil)

type userEntry struct {
	ID          string
	CreatedAt   int64
	UpdatedAt   int64
	Nickname    string
	Password    string
	Email       string
	Verified    bool
	AvatarPath  string
	Description string
}

type User struct {
	logger logging.Logger
	db     database.DB
}

func NewUser(logger logging.Logger, db database.DB) *User {
	return &User{
		logger: logger.With("repository", "sqlite/user"),
		db:     db,
	}
}

func (r *User) Get(ctx context.Context, id model.ID) (model.User, error) {
	const op = "repository.User.Get"

	query := `SELECT * FROM users WHERE id = ?`
	args := []any{id.String()}

	row := r.db.QueryRow(ctx, query, args...)
	user, err := r.scan(row)
	if err != nil {
		if sqlite.IsNoRows(err) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *User) GetByNickname(ctx context.Context, nickname string) (model.User, error) {
	const op = "repository.User.Get"

	query := `SELECT * FROM users WHERE nickname = ?`
	args := []any{nickname}

	row := r.db.QueryRow(ctx, query, args...)
	user, err := r.scan(row)
	if err != nil {
		if sqlite.IsNoRows(err) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *User) GetByEmail(ctx context.Context, email string) (model.User, error) {
	const op = "repository.User.Get"

	query := `SELECT * FROM users WHERE email = ?`
	args := []any{email}

	row := r.db.QueryRow(ctx, query, args...)
	user, err := r.scan(row)
	if err != nil {
		if sqlite.IsNoRows(err) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *User) Create(ctx context.Context, dto repository.CreateUserDTO) (model.ID, error) {
	const op = "repository.User.Create"

	id, err := uuid.NewRandom()
	if err != nil {
		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}
	now := time.Now()

	query := `
		INSERT INTO users (id, created_at, updated_at, nickname, email, password) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	args := []any{id.String(), now.Unix(), now.Unix(), dto.Nickname, dto.Email, dto.Password}

	if err := r.db.Exec(ctx, query, args...); err != nil {
		if sqlite.IsKeyConflict(err) {
			return model.ID{}, fmt.Errorf("%s: %w", op, model.ErrUserExists)
		}

		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return model.ID(id), nil
}

func (r *User) Update(ctx context.Context, id model.ID, dto repository.UpdateUserDTO) error {
	const op = "repository.User.Update"

	now := time.Now()

	query := `UPDATE users SET updated_at = ?`
	args := []any{now.Unix()}

	if dto.Nickname != nil {
		query += `, nickname = ?`
		args = append(args, *dto.Nickname)
	}
	if dto.Password != nil {
		query += `, password = ?`
		args = append(args, *dto.Password)
	}
	if dto.Email != nil {
		query += `, email = ?`
		args = append(args, *dto.Email)
	}
	if dto.Verified != nil {
		query += `, is_verified = ?`
		args = append(args, *dto.Verified)
	}
	if dto.AvatarPath != nil {
		query += `, avatar_path = ?`
		args = append(args, *dto.AvatarPath)
	}
	if dto.Description != nil {
		query += `, description = ?`
		args = append(args, *dto.Description)
	}

	query += ` WHERE id = ?`
	args = append(args, id.String())

	if err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *User) Delete(ctx context.Context, id model.ID) error {
	const op = "repository.User.Delete"

	query := `DELETE FROM users WHERE id = ?`
	args := []any{id.String()}

	if err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (*User) scan(s database.Scanner) (model.User, error) {
	var entry userEntry
	if err := s.Scan(
		&entry.ID, &entry.CreatedAt, &entry.UpdatedAt,
		&entry.Nickname, &entry.Password,
		&entry.Email, &entry.Verified,
		&entry.AvatarPath, &entry.Description,
	); err != nil {
		return model.User{}, err
	}

	id, err := uuid.Parse(entry.ID)
	if err != nil {
		return model.User{}, err
	}

	createdAt := time.Unix(entry.CreatedAt, 0)
	updatedAt := time.Unix(entry.UpdatedAt, 0)

	return model.User{
		Model: model.Model{
			ID:        model.ID(id),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Nickname:    entry.Nickname,
		Email:       entry.Email,
		Password:    entry.Password,
		Verified:    entry.Verified,
		AvatarPath:  entry.AvatarPath,
		Description: entry.Description,
	}, nil
}
