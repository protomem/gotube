package service

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/hashing"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ User = (*UserImpl)(nil)

type CreateUserDTO struct {
	Nickname string
	Password string
	Email    string
}

func (dto CreateUserDTO) Validate() error {
	return nil
}

type (
	User interface {
		Create(ctx context.Context, dto CreateUserDTO) (model.User, error)
	}

	UserImpl struct {
		repo   repository.User
		hasher hashing.Hasher
	}
)

func NewUser(repo repository.User, hasher hashing.Hasher) *UserImpl {
	return &UserImpl{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *UserImpl) Create(ctx context.Context, dto CreateUserDTO) (model.User, error) {
	const op = "service:User.Create"

	if err := dto.Validate(); err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	hash, err := s.hasher.Generate(dto.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	dto.Password = hash

	userID, err := s.repo.Create(ctx, repository.CreateUserDTO(dto))
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.repo.Get(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
