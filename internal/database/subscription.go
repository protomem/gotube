package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var ErrSubscriptionAlreadyExists = NewModelError(ErrAlreadyExists, "subscription")

type Subscription struct {
	ID uuid.UUID `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	FromUserID uuid.UUID `db:"from_user_id" json:"fromUserId"`
	ToUserID   uuid.UUID `db:"to_user_id" json:"toUserId"`
}

func (db *DB) CountSubscriptionsByFromUserID(ctx context.Context, fromUserID uuid.UUID) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT COUNT(*) FROM subscriptions WHERE from_user_id = $1`
	args := []any{fromUserID}

	var count uint64

	if err := db.
		QueryRowxContext(ctx, query, args...).
		Scan(&count); err != nil {
		if IsNoRows(err) {
			return 0, nil

		}

		return 0, err
	}

	return count, nil
}

func (db *DB) CountSubscriptionsByToUserID(ctx context.Context, toUserID uuid.UUID) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT COUNT(*) FROM subscriptions WHERE to_user_id = $1`
	args := []any{toUserID}

	var count uint64

	if err := db.
		QueryRowxContext(ctx, query, args...).
		Scan(&count); err != nil {
		if IsNoRows(err) {
			return 0, nil
		}

		return 0, err
	}

	return count, nil
}

func (db *DB) FindSubscriptionsByFromUserID(ctx context.Context, fromUserID uuid.UUID) ([]Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT * FROM subscriptions WHERE from_user_id = $1`
	args := []any{fromUserID}

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []Subscription{}, nil
		}

		return []Subscription{}, err
	}
	defer func() { _ = rows.Close() }()

	subs := make([]Subscription, 0)
	for rows.Next() {
		var sub Subscription
		if err := rows.StructScan(&sub); err != nil {
			return []Subscription{}, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
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
