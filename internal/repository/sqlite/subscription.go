package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/database/sqlite"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.Subscription = (*Subscription)(nil)

type subscriptionEntry struct {
	ID         string
	CreatedAt  int64
	UpdatedAt  int64
	FromUserID string
	ToUserID   string
}

type Subscription struct {
	logger logging.Logger
	db     database.DB
}

func NewSubscription(logger logging.Logger, db database.DB) *Subscription {
	return &Subscription{
		logger: logger.With("repository", "sqlite/subscription"),
		db:     db,
	}
}

func (r *Subscription) GetByFromUserAndToUser(ctx context.Context, fromUserID, toUserID model.ID) (model.Subscription, error) {
	const op = "repository.Subscription.GetByFromUserAndToUser"

	query := `
		SELECT * FROM subscriptions
		WHERE from_user_id = ? AND to_user_id = ?
		LIMIT 1
	`
	args := []any{fromUserID.String(), toUserID.String()}

	row := r.db.QueryRow(ctx, query, args...)
	sub, err := r.scan(row)
	if err != nil {
		if sqlite.IsNoRows(err) {
			return model.Subscription{}, fmt.Errorf("%s: %w", op, model.ErrSubscriptionNotFound)
		}

		return model.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	r.logger.WithContext(ctx).Debug("found subscription", "subscription", sub)

	return sub, nil
}

func (r *Subscription) CountByToUser(ctx context.Context, toUserID model.ID) (int64, error) {
	const op = "repository.Subscription.CountByToUser"

	query := `SELECT COUNT(id) FROM subscriptions WHERE to_user_id = ?`
	args := []any{toUserID.String()}

	var count int64

	row := r.db.QueryRow(ctx, query, args...)
	if err := row.Scan(&count); err != nil {
		if sqlite.IsNoRows(err) {
			return 0, nil
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (r *Subscription) Create(ctx context.Context, dto repository.CreateSubscriptionDTO) (model.ID, error) {
	const op = "repository.Subscription.Create"

	id, err := uuid.NewRandom()
	if err != nil {
		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}
	now := time.Now()

	query := `
		INSERT INTO subscriptions (id, created_at, updated_at, from_user_id, to_user_id)
		VALUES (?, ?, ?, ?, ?)
	`
	args := []any{id.String(), now.Unix(), now.Unix(), dto.FromUserID.String(), dto.ToUserID.String()}

	if err := r.db.Exec(ctx, query, args...); err != nil {
		if sqlite.IsKeyConflict(err) {
			return model.ID{}, fmt.Errorf("%s: %w", op, model.ErrSubscriptionExists)
		}

		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return model.ID(id), nil
}

func (r *Subscription) Delete(ctx context.Context, id model.ID) error {
	const op = "repository.Subscription.Delete"

	query := `DELETE FROM subscriptions WHERE id = ?`
	args := []any{id.String()}

	if err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Subscription) scan(s database.Scanner) (model.Subscription, error) {
	var entry subscriptionEntry
	if err := s.Scan(
		&entry.ID, &entry.CreatedAt, &entry.UpdatedAt,
		&entry.FromUserID, &entry.ToUserID,
	); err != nil {
		return model.Subscription{}, err
	}

	id, err := uuid.Parse(entry.ID)
	if err != nil {
		return model.Subscription{}, err
	}

	createdAt := time.Unix(entry.CreatedAt, 0)
	updatedAt := time.Unix(entry.UpdatedAt, 0)

	fromUserID, err := uuid.Parse(entry.FromUserID)
	if err != nil {
		return model.Subscription{}, err
	}

	toUserID, err := uuid.Parse(entry.ToUserID)
	if err != nil {
		return model.Subscription{}, err
	}

	return model.Subscription{
		Model: model.Model{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		FromUserID: fromUserID,
		ToUserID:   toUserID,
	}, nil
}
