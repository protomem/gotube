package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/module/user/model"
)

type CreateUserDTO struct {
	Nickname string
	Password string
	Email    string
}

type (
	UserRepository interface {
		FindOneUser(ctx context.Context, id uuid.UUID) (model.User, error)
		CreateUser(ctx context.Context, dto CreateUserDTO) (uuid.UUID, error)
	}
)
