package adapter

import (
	"context"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/infra/database"
)

var _ port.UserAccessor = (*UserAccessor)(nil)

type UserAccessor struct {
	dao *database.UserDAO
}

func NewUserAccessor(db *database.DB) *UserAccessor {
	return &UserAccessor{db.UserDAO()}
}

func (acc *UserAccessor) ByID(ctx context.Context, id entity.ID) (entity.User, error) {
	user, err := acc.dao.GetByID(ctx, id)
	if err != nil {
		if database.IsNoRows(err) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, err
	}

	return entity.User(user), nil
}

var _ port.UserMutator = (*UserMutator)(nil)

type UserMutator struct {
	dao *database.UserDAO
}

func NewUserMutator(db *database.DB) *UserMutator {
	return &UserMutator{db.UserDAO()}
}

func (mut *UserMutator) Insert(ctx context.Context, dto port.InsertUserDTO) (entity.ID, error) {
	id, err := mut.dao.Insert(ctx, database.InsertUserDTO(dto))
	if err != nil {
		if database.IsKeyConflict(err) {
			return entity.ID{}, entity.ErrUserAlreadyExists
		}

		return entity.ID{}, err
	}

	return id, nil
}