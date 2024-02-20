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
	FromUser   userEntry
	ToUserID   string
	ToUser     userEntry
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

func (r *Subscription) GetByFromUserIDAndToUserID(ctx context.Context, fromUserID, toUserID model.ID) (model.Subscription, error) {
	const op = "repository.Subscription.GetByFromUserIDAndToUserID"

	query := `
		SELECT sub.*, from_user.*, to_user.* FROM subscriptions AS sub 
		JOIN users AS from_user ON from_user.id = sub.from_user_id
		JOIN users AS to_user ON to_user.id = sub.to_user_id
		WHERE from_user_id = ? AND to_user_id = ?
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

func (r *Subscription) CountByToUserID(ctx context.Context, toUserID model.ID) (int64, error) {
	const op = "repository.Subscription.CountByToUserID"

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

		&entry.FromUser.ID, &entry.FromUser.CreatedAt, &entry.FromUser.UpdatedAt,
		&entry.FromUser.Nickname, &entry.FromUser.Password,
		&entry.FromUser.Email, &entry.FromUser.Verified,
		&entry.FromUser.AvatarPath, &entry.FromUser.Description,

		&entry.ToUser.ID, &entry.ToUser.CreatedAt, &entry.ToUser.UpdatedAt,
		&entry.ToUser.Nickname, &entry.ToUser.Password,
		&entry.ToUser.Email, &entry.ToUser.Verified,
		&entry.ToUser.AvatarPath, &entry.ToUser.Description,
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

	fromUserCreatedAt := time.Unix(entry.FromUser.CreatedAt, 0)
	fromUserUpdatedAt := time.Unix(entry.FromUser.UpdatedAt, 0)

	toUserID, err := uuid.Parse(entry.ToUserID)
	if err != nil {
		return model.Subscription{}, err
	}

	toUserCreatedAt := time.Unix(entry.ToUser.CreatedAt, 0)
	toUserUpdatedAt := time.Unix(entry.ToUser.UpdatedAt, 0)

	return model.Subscription{
		Model: model.Model{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		FromUser: model.User{
			Model: model.Model{
				ID:        fromUserID,
				CreatedAt: fromUserCreatedAt,
				UpdatedAt: fromUserUpdatedAt,
			},
			Nickname:    entry.FromUser.Nickname,
			Password:    entry.FromUser.Password,
			Email:       entry.FromUser.Email,
			Verified:    entry.FromUser.Verified,
			AvatarPath:  entry.FromUser.AvatarPath,
			Description: entry.FromUser.Description,
		},
		ToUser: model.User{
			Model: model.Model{
				ID:        toUserID,
				CreatedAt: toUserCreatedAt,
				UpdatedAt: toUserUpdatedAt,
			},
			Nickname:    entry.ToUser.Nickname,
			Password:    entry.ToUser.Password,
			Email:       entry.ToUser.Email,
			Verified:    entry.ToUser.Verified,
			AvatarPath:  entry.ToUser.AvatarPath,
			Description: entry.ToUser.Description,
		},
	}, nil
}
