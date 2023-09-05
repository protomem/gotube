package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/module/user/model"
	"github.com/protomem/gotube/internal/module/user/repository"
	"github.com/protomem/gotube/internal/passhash"
)

var _ UserService = (*UserServiceImpl)(nil)

type FindOneUserByEmailAndPasswordDTO struct {
	Email    string
	Password string
}

type CreateUserDTO struct {
	Nickname string
	Password string
	Email    string
}

type (
	UserService interface {
		FindOneUser(ctx context.Context, id uuid.UUID) (model.User, error)
		FindOneUserByEmailAndPassword(ctx context.Context, dto FindOneUserByEmailAndPasswordDTO) (model.User, error)
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

func (s *UserServiceImpl) FindOneUser(ctx context.Context, id uuid.UUID) (model.User, error) {
	const op = "UserService.FindOneUser"

	user, err := s.userRepo.FindOneUser(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserServiceImpl) FindOneUserByEmailAndPassword(
	ctx context.Context,
	dto FindOneUserByEmailAndPasswordDTO,
) (model.User, error) {
	const op = "UserService.FindOneUserByEmailAndPassword"
	var err error

	user, err := s.userRepo.FindOneUserByEmail(ctx, dto.Email)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.hasher.Compare(dto.Password, user.Password)
	if err != nil {
		if errors.Is(err, passhash.ErrWrongPassword) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
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
