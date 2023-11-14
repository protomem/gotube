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

type RatingRepository struct {
	logger logging.Logger
	db     *sql.DB
}

func NewRatingRepository(logger logging.Logger, db *sql.DB) *RatingRepository {
	return &RatingRepository{
		logger: logger.With("repository", "postgres", "repositoryType", "rating"),
		db:     db,
	}
}

func (repo *RatingRepository) FindByVideoID(ctx context.Context, videoID uuid.UUID) ([]model.Rating, error) {
	const op = "repository.Rating.FindByVideoID"

	query := `SELECT * FROM ratings WHERE video_id = $1`

	rows, err := repo.db.QueryContext(ctx, query, videoID)
	if err != nil {
		if pgerr.IsNotFound(err) {
			return []model.Rating{}, nil
		}

		return []model.Rating{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = rows.Close() }()

	var ratings []model.Rating
	for rows.Next() {
		var rating model.Rating
		err := rows.Scan(
			&rating.ID,
			&rating.CreatedAt, &rating.UpdatedAt,
			&rating.Like,
			&rating.VideoID, &rating.UserID,
		)
		if err != nil {
			return ratings, fmt.Errorf("%s: %w", op, err)
		}

		ratings = append(ratings, rating)
	}

	if rows.Err() != nil {
		return ratings, fmt.Errorf("%s: %w", op, rows.Err())
	}

	if len(ratings) == 0 {
		return []model.Rating{}, nil
	}

	return ratings, nil
}

func (repo *RatingRepository) GetByVideoIDAndUserID(
	ctx context.Context,
	videoID, userID uuid.UUID,
) (model.Rating, error) {
	const op = "repository.Rating.GetByVideoIDAndUserID"

	query := `SELECT * FROM ratings WHERE video_id = $1 AND user_id = $2 LIMIT 1`

	var rating model.Rating
	err := repo.db.
		QueryRowContext(ctx, query, videoID, userID).
		Scan(
			&rating.ID,
			&rating.CreatedAt, &rating.UpdatedAt,
			&rating.Like,
			&rating.VideoID, &rating.UserID,
		)
	if err != nil {
		if pgerr.IsNotFound(err) {
			return rating, model.ErrRatingNotFound
		}

		return rating, fmt.Errorf("%s: %w", op, err)
	}

	return rating, nil
}

func (repo *RatingRepository) Create(ctx context.Context, dto repository.CreateRatingDTO) (uuid.UUID, error) {
	const op = "repository.Rating.Create"

	query := `
		INSERT INTO ratings(is_like, video_id, user_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id uuid.UUID
	err := repo.db.
		QueryRowContext(ctx, query, dto.Like, dto.VideoID, dto.UserID).
		Scan(&id)
	if err != nil {
		if pgerr.IsConflict(err) {
			return uuid.Nil, model.ErrRatingExists
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (repo *RatingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.Rating.Delete"

	query := `DELETE FROM ratings WHERE id = $1`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
