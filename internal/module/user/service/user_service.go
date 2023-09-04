package service

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/module/user/model"
	"github.com/protomem/gotube/internal/module/user/repository"
	"github.com/protomem/gotube/internal/passhash"
)

var _ UserService = (*UserServiceImpl)(nil)

type CreateUserDTO struct {
	Nickname string
	Password string
	Email    string
}

type (
	UserService interface {
		CreateUser(ctx context.Context, dto CreateUserDTO) (model.User, error)
	}

	UserServiceImpl struct {
		hasher passhash.Hasher

		userRepo repository.UserRepository
	}
)

func NewUserService(hasher passhash.Hasher, userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		hasher:   hasher,
		userRepo: userRepo,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, dto CreateUserDTO) (model.User, error) {
	const op = "UserService.CreateUser"
	var err error

	// TODO: validate

	dto.Password, err = s.hasher.Generate(dto.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := s.userRepo.CreateUser(ctx, repository.CreateUserDTO(dto))
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userRepo.FindOneUser(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
