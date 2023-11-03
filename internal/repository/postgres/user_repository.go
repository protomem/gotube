package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/pgerr"
)

var _ repository.User = (*UserRepository)(nil)

type UserRepository struct {
	logger logging.Logger
	db     *sql.DB
}

func NewUserRepository(logger logging.Logger, db *sql.DB) *UserRepository {
	return &UserRepository{
		logger: logger.With("repository", "user", "repositoryType", "postgres"),
		db:     db,
	}
}

func (repo *UserRepository) Get(ctx context.Context, id uuid.UUID) (model.User, error) {
	const op = "postgres.UserRepository.Get"

	query := `SELECT * FROM users WHERE id = $1 LIMIT 1`

	var user model.User
	err := repo.db.
		QueryRowContext(ctx, query, id).
		Scan(
			&user.ID,
			&user.CreatedAt, &user.UpdatedAt,
			&user.Nickname, &user.Password,
			&user.Email, &user.Verified,
			&user.AvatarPath, &user.Description,
		)
	if err != nil {
		if pgerr.IsNotFound(err) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (repo *UserRepository) GetByNickname(ctx context.Context, nickname string) (model.User, error) {
	const op = "postgres.UserRepository.GetByNickname"

	query := `SELECT * FROM users WHERE nickname = $1 LIMIT 1`

	var user model.User
	err := repo.db.
		QueryRowContext(ctx, query, nickname).
		Scan(
			&user.ID,
			&user.CreatedAt, &user.UpdatedAt,
			&user.Nickname, &user.Password,
			&user.Email, &user.Verified,
			&user.AvatarPath, &user.Description,
		)
	if err != nil {
		if pgerr.IsNotFound(err) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	const op = "postgres.UserRepository.GetByEmail"

	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`

	var user model.User
	err := repo.db.
		QueryRowContext(ctx, query, email).
		Scan(
			&user.ID,
			&user.CreatedAt, &user.UpdatedAt,
			&user.Nickname, &user.Password,
			&user.Email, &user.Verified,
			&user.AvatarPath, &user.Description,
		)
	if err != nil {
		if pgerr.IsNotFound(err) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (repo *UserRepository) Create(ctx context.Context, dto repository.CreateUserDTO) (uuid.UUID, error) {
	const op = "postgres.UserRepository.Create"

	query := `
        INSERT INTO users (nickname, password, email)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	var id uuid.UUID
	err := repo.db.
		QueryRowContext(ctx, query, dto.Nickname, dto.Password, dto.Email).
		Scan(&id)
	if err != nil {
		if pgerr.IsConflict(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrUserExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (repo *UserRepository) UpdateByNickname(ctx context.Context, nickname string, dto repository.UpdateUserDTO) error {
	const op = "postgres.UserRepository.UpdateByNickname"

	var (
		counter int   = 1
		args    []any = []any{nickname}
		query   strings.Builder
	)
	_, _ = query.WriteString("UPDATE users SET ")

	if dto.Nickname != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("nickname = $%d, ", counter))
		args = append(args, *dto.Nickname)
	}

	if dto.Password != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("password = $%d, ", counter))
		args = append(args, *dto.Password)
	}

	if dto.Email != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("email = $%d, ", counter))
		args = append(args, *dto.Email)
	}

	if dto.Verified != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("is_verified = $%d, ", counter))
		args = append(args, *dto.Verified)
	}

	if dto.AvatarPath != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("avatar_path = $%d, ", counter))
		args = append(args, *dto.AvatarPath)
	}

	if dto.Description != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("description = $%d, ", counter))
		args = append(args, *dto.Description)
	}

	_, _ = query.WriteString("updated_at = now() WHERE nickname = $1")

	_, err := repo.db.ExecContext(ctx, query.String(), args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (repo *UserRepository) DeleteByNickname(ctx context.Context, nickname string) error {
	const op = "postgres.UserRepository.DeleteByNickname"

	query := `DELETE FROM users WHERE nickname = $1`

	_, err := repo.db.ExecContext(ctx, query, nickname)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
