package port

import (
	"context"

	"github.com/protomem/gotube/internal/domain/entity"
)

type InsertUserDTO struct {
	Nickname string
	Password string
	Email    string
}

type (
	UserAccessor interface {
		ByID(ctx context.Context, id entity.ID) (entity.User, error)
	}

	UserMutator interface {
		Insert(ctx context.Context, dto InsertUserDTO) (entity.ID, error)
	}
)
