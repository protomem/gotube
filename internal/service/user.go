package service

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/hashing"
)

var _ User = (*UserImpl)(nil)

type (
	CreateUserDTO struct {
		Nickname string
		Email    string
		Password string
	}
)

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
	const op = "service.User.Create"

	// TODO: add validation

	hashPass, err := s.hasher.Generate(dto.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}
	dto.Password = hashPass

	id, err := s.repo.Create(ctx, repository.CreateUserDTO(dto))
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
