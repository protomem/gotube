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

	UpdateUserDTO struct {
		Nickname    *string
		Password    *string
		Email       *string
		Verified    *bool
		AvatarPath  *string
		Description *string
	}
)

type User interface {
	Get(ctx context.Context, id model.ID) (model.User, error)
	GetByNickname(ctx context.Context, nickname string) (model.User, error)
	Create(ctx context.Context, dto CreateUserDTO) (model.ID, error)
	Update(ctx context.Context, id model.ID, dto UpdateUserDTO) error
}
