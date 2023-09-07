package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/user/model"
	"github.com/protomem/gotube/internal/module/user/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	logger  logging.Logger
	db      *database.DB
	builder squirrel.StatementBuilderType
}

func NewUserRepository(logger logging.Logger, db *database.DB) *UserRepository {
	return &UserRepository{
		logger:  logger.With("repository", "user", "repositoryType", "postgres"),
		db:      db,
		builder: database.Builder(),
	}
}

func (r *UserRepository) FindOneUser(ctx context.Context, id uuid.UUID) (model.User, error) {
	const op = "UserRepository.FindOneUser"
	var err error

	query, args, err := r.builder.
		Select("*").
		From("users").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user model.User
	err = r.db.
		QueryRow(ctx, query, args...).
		Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt,
			&user.Nickname, &user.Password,
			&user.Email, &user.Verified,
			&user.AvatarPath, &user.Description,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserRepository) FindOneUserByNickname(ctx context.Context, nickname string) (model.User, error) {
	const op = "UserRepository.FindOneUserByNickname"
	var err error

	query, args, err := r.builder.
		Select("*").
		From("users").
		Where(squirrel.Eq{"nickname": nickname}).
		Limit(1).
		ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user model.User
	err = r.db.
		QueryRow(ctx, query, args...).
		Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt,
			&user.Nickname, &user.Password,
			&user.Email, &user.Verified,
			&user.AvatarPath, &user.Description,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserRepository) FindOneUserByEmail(ctx context.Context, email string) (model.User, error) {
	const op = "UserRepository.FindOneUserByEmail"
	var err error

	query, args, err := r.builder.
		Select("*").
		From("users").
		Where(squirrel.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user model.User
	err = r.db.
		QueryRow(ctx, query, args...).
		Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt,
			&user.Nickname, &user.Password,
			&user.Email, &user.Verified,
			&user.AvatarPath, &user.Description,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, dto repository.CreateUserDTO) (uuid.UUID, error) {
	const op = "UserRepository.CreateUser"
	var err error

	query, args, err := r.builder.
		Insert("users").
		Columns("nickname", "password", "email").
		Values(dto.Nickname, dto.Password, dto.Email).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	var id uuid.UUID
	err = r.db.
		QueryRow(ctx, query, args...).
		Scan(&id)
	if err != nil {
		if database.IsDuplicateKeyError(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrUserAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *UserRepository) DeleteUserByNickname(ctx context.Context, nickname string) error {
	const op = "UserRepository.DeleteUserByNickname"
	var err error

	query, args, err := r.builder.
		Delete("users").
		Where(squirrel.Eq{"nickname": nickname}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
