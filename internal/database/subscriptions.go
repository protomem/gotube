package database

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/domain/model"
)

func (db *DB) GetSubscriptionByFromUserIDAndToUserID(
	ctx context.Context,
	fromUserID model.ID,
	toUserID model.ID,
) (model.Subscription, error) {
	const op = "database.GetSubscriptionByFromUserIDAndToUserID"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT * FROM subscriptions WHERE from_user_id = $1 AND to_user_id = $2 LIMIT 1`
	args := []any{fromUserID, toUserID}

	var subscription model.Subscription

	if err := db.QueryRowxContext(ctx, query, args...).StructScan(&subscription); err != nil {
		if IsNoRows(err) {
			return model.Subscription{}, fmt.Errorf("%s: %w", op, model.ErrSubscriptionNotFound)
		}

		return model.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	return subscription, nil
}

type InsertSubscriptionDTO struct {
	FromUserID model.ID
	ToUserID   model.ID
}

func (db *DB) InsertSubscription(ctx context.Context, dto InsertSubscriptionDTO) (model.ID, error) {
	const op = "database.InsertSubscription"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `INSERT INTO subscriptions (from_user_id, to_user_id) VALUES ($1, $2) RETURNING id`
	args := []any{dto.FromUserID, dto.ToUserID}

	var id model.ID

	if err := db.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		if IsKeyConflict(err) {
			return model.ID{}, fmt.Errorf("%s: %w", op, model.ErrSubscriptionAlreadyExists)
		}

		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (db *DB) DeleteSubscription(ctx context.Context, id model.ID) error {
	const op = "database.DeleteSubscription"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `DELETE FROM subscriptions WHERE id = $1`
	args := []any{id}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
