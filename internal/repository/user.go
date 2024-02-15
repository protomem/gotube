package repository

import (
	"context"

	"github.com/protomem/gotube/internal/model"
)

type (
	CreateUserDTO struct {
		Nickname string
		Email    string
		Password string
	}
)

type User interface {
	Get(ctx context.Context, id model.ID) (model.User, error)
	GetByNickname(ctx context.Context, nickname string) (model.User, error)
	Create(ctx context.Context, dto CreateUserDTO) (model.ID, error)
}
