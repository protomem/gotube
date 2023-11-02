package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ repository.User = (*UserRepository)(nil)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) Get(ctx context.Context, id uuid.UUID) (model.User, error) {
	panic("not implemented")
}

func (repo *UserRepository) GetByNickname(ctx context.Context, nickname string) (model.User, error) {
	panic("not implemented")
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	panic("not implemented")
}

func (repo *UserRepository) Create(ctx context.Context, dto repository.CreateUserDTO) error {
	panic("not implemented")
}

func (repo *UserRepository) UpdateByNickname(ctx context.Context, nickname string, dto repository.UpdateUserDTO) error {
	panic("not implemented")
}

func (repo *UserRepository) DeleteByNickname(ctx context.Context, nickname string) error {
	panic("not implemented")
}
