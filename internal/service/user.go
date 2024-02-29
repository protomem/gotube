package service

import (
	"context"
	"errors"
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

	UpdateUserDTO struct {
		Nickname    *string
		Email       *string
		Verified    *bool
		AvatarPath  *string
		Description *string

		NewPassword *string
		OldPassword *string
	}
)

type (
	User interface {
		GetByNickname(ctx context.Context, nickname string) (model.User, error)
		GetByEmailAndPassword(ctx context.Context, email, password string) (model.User, error)
		Create(ctx context.Context, dto CreateUserDTO) (model.User, error)
		UpdateByNickname(ctx context.Context, nickname string, dto UpdateUserDTO) (model.User, error)
		DeleteByNickname(ctx context.Context, nickname string) error
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

func (s *UserImpl) GetByNickname(ctx context.Context, nickname string) (model.User, error) {
	const op = "service.User.GetByNickname"

	user, err := s.repo.GetByNickname(ctx, nickname)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserImpl) GetByEmailAndPassword(ctx context.Context, email, password string) (model.User, error) {
	const op = "service.User.GetByEmailAndPassword"

	// TODO: add validation

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.hasher.Verify(password, user.Password); err != nil {
		if errors.Is(err, hashing.ErrWrongPassword) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
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

func (s *UserImpl) UpdateByNickname(ctx context.Context, nickname string, dto UpdateUserDTO) (model.User, error) {
	const op = "service.User.UpdateByNickname"

	// TODO: add validation

	oldUser, err := s.repo.GetByNickname(ctx, nickname)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	repoDTO := repository.UpdateUserDTO{
		Nickname:    dto.Nickname,
		Email:       dto.Email,
		Verified:    dto.Verified,
		AvatarPath:  dto.AvatarPath,
		Description: dto.Description,
	}

	if dto.Email != nil {
		repoDTO.Verified = new(bool)
		*repoDTO.Verified = false
	}

	if dto.NewPassword != nil && dto.OldPassword != nil {
		if err := s.hasher.Verify(*dto.OldPassword, oldUser.Password); err != nil {
			if errors.Is(err, hashing.ErrWrongPassword) {
				return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
			}

			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}

		newHashPass, err := s.hasher.Generate(*dto.NewPassword)
		if err != nil {
			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}
		repoDTO.Password = &newHashPass
	}

	if err := s.repo.Update(ctx, oldUser.ID, repoDTO); err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	newUser, err := s.repo.Get(ctx, oldUser.ID)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return newUser, nil
}

func (s *UserImpl) DeleteByNickname(ctx context.Context, nickname string) error {
	const op = "service.User.DeleteByNickname"

	user, err := s.repo.GetByNickname(ctx, nickname)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.repo.Delete(ctx, user.ID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
