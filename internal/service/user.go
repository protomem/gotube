package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/hashing"
	"github.com/protomem/gotube/internal/hashing/bcrypt"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ User = (*UserImpl)(nil)

var ErrInvalidPassword = errors.New("invalid password")

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

		UpdateByNickname(ctx context.Context, nickname string, dto UpdateUserDTO) (model.User, error)

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
	const op = "service.User.Get"

	user, err := serv.repo.Get(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (serv *UserImpl) GetByNickname(ctx context.Context, nickname string) (model.User, error) {
	const op = "service.User.GetByNickname"

	// TODO: Valiate ...

	user, err := serv.repo.GetByNickname(ctx, nickname)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (serv *UserImpl) GetByEmailAndPassword(ctx context.Context, email, password string) (model.User, error) {
	const op = "service.User.GetByEmailAndPassword"
	var err error

	user, err := serv.repo.GetByEmail(ctx, email)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	err = serv.hasher.Compare(password, user.Password)
	if err != nil {
		if errors.Is(err, hashing.ErrWrongPassword) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (serv *UserImpl) Create(ctx context.Context, dto CreateUserDTO) (model.User, error) {
	const op = "service.User.Create"
	var err error

	// TODO: Valiate ...

	passHash, err := serv.hasher.Generate(dto.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	dto.Password = passHash

	id, err := serv.repo.Create(ctx, repository.CreateUserDTO(dto))
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := serv.repo.Get(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (serv *UserImpl) UpdateByNickname(ctx context.Context, nickname string, dto UpdateUserDTO) (model.User, error) {
	const op = "service.User.UpdateByNickname"
	var err error

	// TODO: Validate ...

	oldUser, err := serv.repo.GetByNickname(ctx, *dto.Nickname)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	updateData := repository.UpdateUserDTO{}

	if dto.Nickname != nil {
		updateData.Nickname = dto.Nickname
	}

	if dto.Email != nil {
		updateData.Email = dto.Email
	}

	if dto.AvatarPath != nil {
		updateData.AvatarPath = dto.AvatarPath
	}

	if dto.Description != nil {
		updateData.Description = dto.AvatarPath
	}

	if dto.NewPassword != nil && dto.OldPassword != nil {
		err = serv.hasher.Compare(*dto.OldPassword, oldUser.Password)
		if err != nil {
			if errors.Is(err, hashing.ErrWrongPassword) {
				return model.User{}, fmt.Errorf("%s: %w", op, err)
			}

			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}

		newPassHash, err := serv.hasher.Generate(*dto.NewPassword)
		if err != nil {
			return model.User{}, fmt.Errorf("%s: %w", op, err)
		}

		updateData.Password = &newPassHash
	}

	err = serv.repo.UpdateByNickname(ctx, nickname, updateData)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	newUser, err := serv.repo.Get(ctx, oldUser.ID)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return newUser, nil
}

func (serv *UserImpl) DeleteByNickname(ctx context.Context, nickname string) error {
	const op = "service.User.DeleteByNickname"

	// TODO: Validate ...

	err := serv.repo.DeleteByNickname(ctx, nickname)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
