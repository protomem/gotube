package service

import (
	"context"
	"errors"
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
		FindByIDs(ctx context.Context, ids ...model.ID) ([]model.User, error)
		Get(ctx context.Context, id model.ID) (model.User, error)
		GetByNickname(ctx context.Context, nickname string) (model.User, error)
		GetByEmailAndPassword(ctx context.Context, email, password string) (model.User, error)
		Create(ctx context.Context, dto CreateUserDTO) (model.User, error)
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

func (s *UserImpl) FindByIDs(ctx context.Context, ids ...model.ID) ([]model.User, error) {
	const op = "service:User.FindByIDs"

	users, err := s.repo.FindByIDs(ctx, ids...)
	if err != nil {
		return []model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (s *UserImpl) Get(ctx context.Context, id model.ID) (model.User, error) {
	const op = "service:User.Get"

	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserImpl) GetByNickname(ctx context.Context, nickname string) (model.User, error) {
	const op = "service:User.GetByNickname"

	user, err := s.repo.GetByNickname(ctx, nickname)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserImpl) GetByEmailAndPassword(ctx context.Context, email, password string) (model.User, error) {
	const op = "service:User.GetByEmailAndPassword"

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.hasher.Compare(password, user.Password); err != nil {
		if errors.Is(err, hashing.ErrWrongPassword) {
			return model.User{}, model.ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
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

func (s *UserImpl) DeleteByNickname(ctx context.Context, nickname string) error {
	const op = "service:User.DeleteByNickname"

	if err := s.repo.DeleteByNickname(ctx, nickname); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
