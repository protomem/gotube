package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/hashing"
	"github.com/protomem/gotube/internal/hashing/bcrypt"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ User = (*UserImpl)(nil)

type (
	CreateUserDTO struct {
		Nickname string
		Password string
		Email    string
	}

	UpdateUserDTO struct {
		Nickname    *string
		Email       *string
		AvatarPath  *string
		Description *string

		OldPassword *string
		NewPassword *string
	}

	User interface {
		Get(ctx context.Context, id uuid.UUID) (model.User, error)
		GetByNickname(ctx context.Context, nickname string) (model.User, error)
		GetByEmailAndPassword(ctx context.Context, email, password string) (model.User, error)

		Create(ctx context.Context, dto CreateUserDTO) (model.User, error)

		UpdateByNickname(ctx context.Context, nickname string, dto UpdateUserDTO) error

		DeleteByNickname(ctx context.Context, nickname string) error
	}

	UserImpl struct {
		repo   repository.User
		hasher hashing.Hasher
	}
)

func NewUser(repo repository.User) *UserImpl {
	return &UserImpl{
		repo:   repo,
		hasher: bcrypt.New(bcrypt.DefaultCost),
	}
}

func (serv *UserImpl) Get(ctx context.Context, id uuid.UUID) (model.User, error) {
	panic("not implemented")
}

func (serv *UserImpl) GetByNickname(ctx context.Context, nickname string) (model.User, error) {
	panic("not implemented")
}

func (serv *UserImpl) GetByEmailAndPassword(ctx context.Context, email, password string) (model.User, error) {
	panic("not implemented")
}

func (serv *UserImpl) Create(ctx context.Context, dto CreateUserDTO) (model.User, error) {
	panic("not implemented")
}

func (serv *UserImpl) UpdateByNickname(ctx context.Context, nickname string, dto UpdateUserDTO) error {
	panic("not implemented")
}

func (serv *UserImpl) DeleteByNickname(ctx context.Context, nickname string) error {
	panic("not implemented")
}
