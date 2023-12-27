package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSubscriptionNotFound      = NewModelError(ErrNotFound, "subscription")
	ErrSubscriptionAlreadyExists = NewModelError(ErrAlreadyExists, "subscription")
)

type Subscription struct {
	ID uuid.UUID `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	FromUserID uuid.UUID `db:"from_user_id" json:"fromUserId"`
	ToUserID   uuid.UUID `db:"to_user_id" json:"toUserId"`
}

type CreateSubscriptionDTO struct {
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
}

func (db *DB) CreateSubscription(ctx context.Context, dto CreateSubscriptionDTO) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `INSERT INTO subscriptions(from_user_id, to_user_id) VALUES ($1, $2) RETURNING id`
	args := []any{dto.FromUserID, dto.ToUserID}

	var id uuid.UUID

	if err := db.
		QueryRowxContext(ctx, query, args...).
		Scan(&id); err != nil {
		if IsKeyConflict(err) {
			return uuid.Nil, ErrSubscriptionAlreadyExists
		}

		return uuid.Nil, err
	}

	return id, nil
}

type DeleteSubscriptionDTO struct {
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
}

func (db *DB) DeleteSubscription(ctx context.Context, dto DeleteSubscriptionDTO) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `DELETE FROM subscriptions WHERE from_user_id = $1 AND to_user_id = $2`
	args := []any{dto.FromUserID, dto.ToUserID}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
