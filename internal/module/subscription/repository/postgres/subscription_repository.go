package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/subscription/model"
	"github.com/protomem/gotube/internal/module/subscription/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.SubscriptionRepository = (*SubscriptionRepository)(nil)

type SubscriptionRepository struct {
	logger  logging.Logger
	db      *database.DB
	builder squirrel.StatementBuilderType
}

func NewSubscriptionRepository(logger logging.Logger, db *database.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		logger:  logger.With("repository", "subscription", "repositoryType", "postgres"),
		db:      db,
		builder: database.Builder(),
	}
}

func (r *SubscriptionRepository) CreateSubscription(
	ctx context.Context,
	dto repository.CreateSubscriptionDTO,
) (uuid.UUID, error) {
	const op = "SubscriptionRepository.CreateSubscription"
	var err error

	query, args, err := r.builder.
		Insert("subscriptions").
		Columns("from_user_id", "to_user_id").
		Values(dto.FromUserID, dto.ToUserID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	var id uuid.UUID
	err = r.db.
		QueryRow(ctx, query, args...).
		Scan(&id)
	if err != nil {
		if database.IsDuplicateKeyError(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrSubscriptionAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
