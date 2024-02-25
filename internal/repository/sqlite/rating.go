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

var _ repository.Rating = (*Rating)(nil)

type ratingEntry struct {
	ID        string
	CreatedAt int64
	UpdatedAt int64
	VideoID   string
	UserID    string
	Like      bool
}

type Rating struct {
	logger logging.Logger
	db     database.DB
}

func NewRating(logger logging.Logger, db database.DB) *Rating {
	return &Rating{
		logger: logger.With("repository", "sqlite/rating"),
		db:     db,
	}
}

func (r *Rating) Get(ctx context.Context, dto repository.RatingDTO) (model.Rating, error) {
	const op = "repository.Rating.Get"

	query := `SELECT * FROM ratings WHERE video_id = ? AND user_id = ? LIMIT 1`
	args := []any{dto.VideoID.String(), dto.UserID.String()}

	row := r.db.QueryRow(ctx, query, args...)
	rating, err := r.scan(row)
	if err != nil {
		if sqlite.IsNoRows(err) {
			return model.Rating{}, fmt.Errorf("%s: %w", op, model.ErrRatingNotFound)
		}

		return model.Rating{}, fmt.Errorf("%s: %w", op, err)
	}

	return rating, nil
}

func (r *Rating) Create(ctx context.Context, dto repository.CreateRatingDTO) (model.ID, error) {
	const op = "repository.Rating.Create"

	id, err := uuid.NewRandom()
	if err != nil {
		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}
	now := time.Now()

	query := `INSERT INTO ratings (id, created_at, updated_at, video_id, user_id, is_like) VALUES (?, ?, ?, ?, ?, ?)`
	args := []any{id.String(), now.Unix(), now.Unix(), dto.VideoID.String(), dto.UserID.String(), dto.Like}

	if err := r.db.Exec(ctx, query, args...); err != nil {
		if sqlite.IsKeyConflict(err) {
			return model.ID{}, fmt.Errorf("%s: %w", op, model.ErrRatingExists)
		}

		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Rating) Delete(ctx context.Context, dto repository.RatingDTO) error {
	const op = "repository.Rating.Delete"

	query := `DELETE FROM ratings WHERE video_id = ? AND user_id = ?`
	args := []any{dto.VideoID.String(), dto.UserID.String()}

	if err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Rating) scan(s database.Scanner) (model.Rating, error) {
	var entry ratingEntry
	if err := s.Scan(
		&entry.ID, &entry.CreatedAt, &entry.UpdatedAt,
		&entry.UserID, &entry.VideoID,
		&entry.Like,
	); err != nil {
		return model.Rating{}, err
	}

	id, err := uuid.Parse(entry.ID)
	if err != nil {
		return model.Rating{}, err
	}

	createdAt := time.Unix(entry.CreatedAt, 0)
	updatedAt := time.Unix(entry.UpdatedAt, 0)

	userID, err := uuid.Parse(entry.UserID)
	if err != nil {
		return model.Rating{}, err
	}

	videoID, err := uuid.Parse(entry.VideoID)
	if err != nil {
		return model.Rating{}, err
	}

	return model.Rating{
		Model: model.Model{
			ID:        model.ID(id),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		UserID:  model.ID(userID),
		VideoID: model.ID(videoID),
		Like:    entry.Like,
	}, nil
}
