package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/database/postgres"
	"github.com/protomem/gotube/internal/module/video/model"
	"github.com/protomem/gotube/internal/module/video/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.RatingRepository = (*RatingRepository)(nil)

type RatingRepository struct {
	logger  logging.Logger
	db      *postgres.DB
	builder squirrel.StatementBuilderType
}

func NewRatingRepository(logger logging.Logger, db *postgres.DB) *RatingRepository {
	return &RatingRepository{
		logger:  logger.With("repository", "rating", "repositoryType", "postgres"),
		db:      db,
		builder: postgres.Builder(),
	}
}

func (r *RatingRepository) FindAllRatingsByVideoID(ctx context.Context, videoID uuid.UUID) ([]model.Rating, error) {
	const op = "RatingRepository.FindAllRatingsByVideoID"
	var err error

	qeury, args, err := r.builder.
		Select("*").
		From("video_ratings").
		Where(squirrel.Eq{"video_id": videoID}).
		ToSql()
	if err != nil {
		return []model.Rating{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.db.Query(ctx, qeury, args...)
	if err != nil {
		return []model.Rating{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var ratings []model.Rating
	for rows.Next() {
		var rating model.Rating
		err = rows.Scan(
			&rating.ID, &rating.CreatedAt, &rating.UpdatedAt,
			&rating.VideoID, &rating.UserID,
			&rating.Type,
		)
		if err != nil {
			return []model.Rating{}, fmt.Errorf("%s: %w", op, err)
		}

		ratings = append(ratings, rating)
	}

	if len(ratings) == 0 {
		return []model.Rating{}, nil
	}

	return ratings, nil
}

func (r *RatingRepository) CreateRating(ctx context.Context, dto repository.CreateRatingDTO) (uuid.UUID, error) {
	const op = "RatingRepository.CreateRating"
	var err error

	query, args, err := r.builder.
		Insert("video_ratings").
		Columns("user_id", "video_id", "rating").
		Values(dto.UserID, dto.VideoID, dto.Type).
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
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrRatingAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *RatingRepository) DeleteRating(ctx context.Context, dto repository.DeleteRatingDTO) error {
	const op = "RatingRepository.DeleteRating"
	var err error

	query, args, err := s.builder.
		Delete("video_ratings").
		Where(squirrel.Eq{"user_id": dto.UserID, "video_id": dto.VideoID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
