package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/video/model"
	"github.com/protomem/gotube/internal/module/video/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.RatingRepository = (*RatingRepository)(nil)

type RatingRepository struct {
	logger  logging.Logger
	db      *database.DB
	builder squirrel.StatementBuilderType
}

func NewRatingRepository(logger logging.Logger, db *database.DB) *RatingRepository {
	return &RatingRepository{
		logger:  logger.With("repository", "rating", "repositoryType", "postgres"),
		db:      db,
		builder: database.Builder(),
	}
}

func (r *RatingRepository) CreateRating(ctx context.Context, dto repository.CreateRatingDTO) (uuid.UUID, error) {
	const op = "RatingRepository.CreateRating"
	var err error

	query, args, err := r.builder.
		Insert("video_ratings").
		Columns("user_id", "video_id", "rating").
		Values(dto.UserID, dto.VideID, dto.Type).
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
		if database.IsDuplicateKeyError(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrRatingAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
