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

type VideoRepository struct {
	logger logging.Logger
	db     *sql.DB
}

func NewVideoRepository(logger logging.Logger, db *sql.DB) *VideoRepository {
	return &VideoRepository{
		logger: logger.With("repository", "video", "repositoryType", "postgres"),
		db:     db,
	}
}

func (repo *VideoRepository) Get(ctx context.Context, id uuid.UUID) (model.Video, error) {
	const op = "postgres.VideoRepository.Get"

	query := `
		SELECT videos.*, authors.*
		FROM videos
		JOIN users AS authors ON videos.author_id = authors.id
		WHERE videos.id = $1
		LIMIT 1
	`

	var video model.Video
	err := repo.db.
		QueryRowContext(ctx, query, id).
		Scan(
			&video.ID,
			&video.CreatedAt, &video.UpdatedAt,
			&video.Title, &video.Description,
			&video.ThumbnailPath, &video.VideoPath,
			&video.Author.ID, &video.Public, &video.Views,

			&video.Author.ID,
			&video.Author.CreatedAt, &video.Author.UpdatedAt,
			&video.Author.Nickname, &video.Author.Password,
			&video.Author.Email, &video.Author.Verified,
			&video.Author.AvatarPath, &video.Author.Description,
		)
	if err != nil {
		if pgerr.IsNotFound(err) {
			return model.Video{}, fmt.Errorf("%s: %w", op, model.ErrVideoNotFound)
		}

		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (repo *VideoRepository) Create(ctx context.Context, dto repository.CreateVideoDTO) (uuid.UUID, error) {
	const op = "postgres.VideoRepository.Create"

	query := `
		INSERT INTO videos (title, description, thumbnail_path, video_path, author_id, is_public)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id uuid.UUID
	err := repo.db.
		QueryRowContext(
			ctx,
			query,
			dto.Title, dto.Description,
			dto.ThumbnailPath, dto.VideoPath,
			dto.AuthorID, dto.Public,
		).
		Scan(&id)
	if err != nil {
		if pgerr.IsConflict(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrVideoExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
