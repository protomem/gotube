package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/protomem/gotube/internal/database/postgres"
	"github.com/protomem/gotube/internal/module/subscription/model"
	"github.com/protomem/gotube/internal/module/subscription/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.SubscriptionRepository = (*SubscriptionRepository)(nil)

type SubscriptionRepository struct {
	logger  logging.Logger
	db      *postgres.DB
	builder squirrel.StatementBuilderType
}

func NewSubscriptionRepository(logger logging.Logger, db *postgres.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		logger:  logger.With("repository", "subscription", "repositoryType", "postgres"),
		db:      db,
		builder: postgres.Builder(),
	}
}

func (r *SubscriptionRepository) FindAllSubscriptionsByFromUserID(
	ctx context.Context,
	fromUserID uuid.UUID,
) ([]model.Subscription, error) {
	const op = "SubscriptionRepository.FindAllSubscriptionsByFromUserID"
	var err error

	query, args, err := r.builder.
		Select("subscriptions.*, from_users.*, to_users.*").
		From("subscriptions").
		Where(squirrel.Eq{
			"from_user_id": fromUserID,
		}).
		Join("users AS from_users ON from_users.id = subscriptions.from_user_id").
		Join("users AS to_users ON to_users.id = subscriptions.to_user_id").
		ToSql()
	if err != nil {
		return []model.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.Subscription{}, nil
		}

		return []model.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	subscriptions := []model.Subscription{}
	for rows.Next() {
		var subscription model.Subscription
		err = rows.Scan(
			&subscription.ID, &subscription.CreatedAt, &subscription.UpdatedAt,

			&subscription.FromUser.ID, &subscription.ToUser.ID,

			&subscription.FromUser.ID, &subscription.FromUser.CreatedAt, &subscription.UpdatedAt,
			&subscription.FromUser.Nickname, &subscription.FromUser.Password,
			&subscription.FromUser.Email, &subscription.FromUser.Verified,
			&subscription.FromUser.AvatarPath, &subscription.FromUser.Description,

			&subscription.ToUser.ID, &subscription.ToUser.CreatedAt, &subscription.ToUser.UpdatedAt,
			&subscription.ToUser.Nickname, &subscription.ToUser.Password,
			&subscription.ToUser.Email, &subscription.ToUser.Verified,
			&subscription.ToUser.AvatarPath, &subscription.ToUser.Description,
		)
		if err != nil {
			return []model.Subscription{}, fmt.Errorf("%s: %w", op, err)
		}

		r.logger.Debug("find all subscriptions by rows", "subscription", subscription)

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
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
		if postgres.IsDuplicateKeyError(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrSubscriptionAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *SubscriptionRepository) DeleteSubscription(
	ctx context.Context,
	dto repository.DeleteSubscriptionDTO,
) error {
	const op = "SubscriptionRepository.DeleteSubscription"
	var err error

	query, args, err := r.builder.
		Delete("subscriptions").
		Where(squirrel.Eq{
			"from_user_id": dto.FromUserID,
			"to_user_id":   dto.ToUserID,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
