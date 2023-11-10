package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/pgerr"
)

var _ repository.Subscription = (*SubscriptionRepository)(nil)

type SubscriptionRepository struct {
	logger logging.Logger
	db     *sql.DB
}

func NewSubscriptionRepository(logger logging.Logger, db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		logger: logger.With("repository", "subscription", "repositoryType", "postgres"),
		db:     db,
	}
}

func (repo *SubscriptionRepository) FindByFromUserID(
	ctx context.Context,
	fromUserID uuid.UUID,
) ([]model.Subscription, error) {
	const op = "postgres.SubscriptionRepository.FindByFromUserID"

	query := `
		SELECT subscriptions.*, from_users.*, to_users.*
		FROM subscriptions
		JOIN users as from_users ON subscriptions.from_user_id = from_users.id
		JOIN users as to_users ON subscriptions.to_user_id = to_users.id
		WHERE from_user_id = $1
	`

	rows, err := repo.db.QueryContext(ctx, query, fromUserID)
	if err != nil {
		if pgerr.IsNotFound(err) {
			return []model.Subscription{}, nil
		}

		return []model.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	var subs []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.CreatedAt, &sub.UpdatedAt,
			&sub.FromUser.ID, &sub.ToUser.ID,

			&sub.FromUser.ID,
			&sub.FromUser.CreatedAt, &sub.FromUser.UpdatedAt,
			&sub.FromUser.Nickname, &sub.FromUser.Password,
			&sub.FromUser.Email, &sub.FromUser.Verified,
			&sub.FromUser.AvatarPath, &sub.FromUser.Description,

			&sub.ToUser.ID,
			&sub.ToUser.CreatedAt, &sub.ToUser.UpdatedAt,
			&sub.ToUser.Nickname, &sub.ToUser.Password,
			&sub.ToUser.Email, &sub.ToUser.Verified,
			&sub.ToUser.AvatarPath, &sub.ToUser.Description,
		)
		if err != nil {
			return []model.Subscription{}, fmt.Errorf("%s: %w", op, err)
		}

		subs = append(subs, sub)
	}

	if rows.Err() != nil {
		return []model.Subscription{}, fmt.Errorf("%s: %w", op, rows.Err())
	}

	if len(subs) == 0 {
		return []model.Subscription{}, nil
	}

	return subs, nil
}

func (repo *SubscriptionRepository) GetByFromUserAndToUser(
	ctx context.Context,
	fromUserID, toUserID uuid.UUID,
) (model.Subscription, error) {
	const op = "postgres.SubscriptionRepository.GetByFromUserAndToUser"

	query := `
        SELECT subscriptions.*, from_users.*, to_users.*
        FROM subscriptions
        JOIN users as from_users ON subscriptions.from_user_id = from_users.id
        JOIN users as to_users ON subscriptions.to_user_id = to_users.id
        WHERE from_user_id = $1 AND to_user_id = $2 
        LIMIT 1
    `

	var sub model.Subscription
	err := repo.db.
		QueryRowContext(ctx, query, fromUserID, toUserID).
		Scan(
			&sub.ID,
			&sub.CreatedAt, &sub.UpdatedAt,
			&sub.FromUser.ID, &sub.ToUser.ID,

			&sub.FromUser.ID,
			&sub.FromUser.CreatedAt, &sub.FromUser.UpdatedAt,
			&sub.FromUser.Nickname, &sub.FromUser.Password,
			&sub.FromUser.Email, &sub.FromUser.Verified,
			&sub.FromUser.AvatarPath, &sub.FromUser.Description,

			&sub.ToUser.ID,
			&sub.ToUser.CreatedAt, &sub.ToUser.UpdatedAt,
			&sub.ToUser.Nickname, &sub.ToUser.Password,
			&sub.ToUser.Email, &sub.ToUser.Verified,
			&sub.ToUser.AvatarPath, &sub.ToUser.Description,
		)
	if err != nil {
		if pgerr.IsNotFound(err) {
			return model.Subscription{}, fmt.Errorf("%s: %w", op, model.ErrSubscriptionNotFound)
		}

		return model.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}

func (repo *SubscriptionRepository) Create(
	ctx context.Context,
	dto repository.CreateSubscriptionDTO,
) (uuid.UUID, error) {
	const op = "postgres.SubscriptionRepository.Create"

	query := `
        INSERT INTO subscriptions (from_user_id, to_user_id)
        VALUES ($1, $2)
        RETURNING id
    `

	var id uuid.UUID
	err := repo.db.
		QueryRowContext(ctx, query, dto.FromUserID, dto.ToUserID).
		Scan(&id)
	if err != nil {
		if pgerr.IsConflict(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrSubscriptionExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (repo *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "postgres.SubscriptionRepository.Delete"

	query := `DELETE FROM subscriptions WHERE id = $1`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
