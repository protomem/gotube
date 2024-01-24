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
		ByNickname(ctx context.Context, nickname string) (entity.User, error)
		ByEmail(ctx context.Context, email string) (entity.User, error)
	}

	UserMutator interface {
		Insert(ctx context.Context, dto InsertUserDTO) (entity.ID, error)
		Delete(ctx context.Context, id entity.ID) error
	}
)
