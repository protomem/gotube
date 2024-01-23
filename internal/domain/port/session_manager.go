package port

import (
	"context"

	"github.com/protomem/gotube/internal/domain/entity"
)

type SessionManager interface {
	Get(ctx context.Context, token string) (entity.Session, error)
	Put(ctx context.Context, session entity.Session) error
	Delete(ctx context.Context, token string) error
}
