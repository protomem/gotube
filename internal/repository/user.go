package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/pkg/logging"
)

var _ User = (*UserImpl)(nil)

type CreateUserDTO struct {
	Nickname string
	Password string
	Email    string
}

type (
	User interface {
		Get(ctx context.Context, id model.ID) (model.User, error)
		GetByNickname(ctx context.Context, nickname string) (model.User, error)
		Create(ctx context.Context, dto CreateUserDTO) (model.ID, error)
	}

	UserImpl struct {
		logger logging.Logger
		pdb    *pgxpool.Pool
	}
)

func NewUser(logger logging.Logger, pdb *pgxpool.Pool) *UserImpl {
	return &UserImpl{
		logger: logger.With("repository", "user", "repositoryType", "postgres"),
		pdb:    pdb,
	}
}

func (r *UserImpl) Get(ctx context.Context, id model.ID) (model.User, error) {
	const op = "repository:User.Get"

	query := `SELECT * FROM users WHERE id = $1 LIMIT 1`
	args := []any{id}

	row := r.pdb.QueryRow(ctx, query, args...)

	var user model.User
	if err := r.scan(row, &user); err != nil {
		if IsPgNotFound(err) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserImpl) GetByNickname(ctx context.Context, nickname string) (model.User, error) {
	const op = "repository:User.GetByNickname"

	query := `SELECT * FROM users WHERE nickname = $1 LIMIT 1`
	args := []any{nickname}

	row := r.pdb.QueryRow(ctx, query, args...)

	var user model.User
	if err := r.scan(row, &user); err != nil {
		if IsPgNotFound(err) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserImpl) Create(ctx context.Context, dto CreateUserDTO) (model.ID, error) {
	const op = "repository:User.Create"

	query := `INSERT INTO users (nickname, password, email) VALUES ($1, $2, $3) RETURNING id`
	args := []any{dto.Nickname, dto.Password, dto.Email}

	row := r.pdb.QueryRow(ctx, query, args...)

	var id model.ID
	if err := row.Scan(&id); err != nil {
		if IsPgDuplicateKey(err) {
			return model.ID{}, fmt.Errorf("%s: %w", op, model.ErrUserAlreadyExists)
		}

		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *UserImpl) scan(row pgx.Row, user *model.User) error {
	return row.Scan(
		&user.ID,
		&user.CreatedAt, &user.UpdatedAt,
		&user.Nickname, &user.Password,
		&user.Email, &user.Verified,
		&user.AvatarPath, &user.Description,
	)
}
