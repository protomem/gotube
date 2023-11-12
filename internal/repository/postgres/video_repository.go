package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
	"github.com/protomem/gotube/pkg/pgerr"
)

var _ repository.Video = (*VideoRepository)(nil)

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

func (repo *VideoRepository) FindAllPublicSortByCreatedAt(
	ctx context.Context,
	opts repository.FindVideosOptions,
) ([]model.Video, error) {
	const op = "postgres.VideoRepository.FindAllPublicSortByCreatedAt"

	query := `
		SELECT videos.*, authors.*
		FROM videos
		JOIN users AS authors ON videos.author_id = authors.id
		WHERE videos.is_public = true
		ORDER BY videos.created_at DESC
		LIMIT $1
		OFFSET $2
	`

	rows, err := repo.db.QueryContext(ctx, query, opts.Limit, opts.Offset)
	if err != nil {
		if pgerr.IsNotFound(err) {
			return []model.Video{}, nil
		}

		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = rows.Close() }()

	var videos []model.Video
	for rows.Next() {
		var video model.Video
		err := rows.Scan(
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
			return []model.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		videos = append(videos, video)
	}

	if rows.Err() != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, rows.Err())
	}

	if len(videos) == 0 {
		return []model.Video{}, nil
	}

	return videos, nil
}

func (repo *VideoRepository) FindAllPublicSortByViews(
	ctx context.Context,
	opts repository.FindVideosOptions,
) ([]model.Video, error) {
	const op = "postgres.VideoRepository.FindAllPublicSortByViews"

	query := `
		SELECT videos.*, authors.*
		FROM videos
		JOIN users AS authors ON videos.author_id = authors.id
		WHERE videos.is_public = true
		ORDER BY videos.views DESC
		LIMIT $1
		OFFSET $2
	`

	rows, err := repo.db.QueryContext(ctx, query, opts.Limit, opts.Offset)
	if err != nil {
		if pgerr.IsNotFound(err) {
			return []model.Video{}, nil
		}

		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = rows.Close() }()

	var videos []model.Video
	for rows.Next() {
		var video model.Video
		err := rows.Scan(
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
			return []model.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		videos = append(videos, video)
	}

	if rows.Err() != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, rows.Err())
	}

	if len(videos) == 0 {
		return []model.Video{}, nil
	}

	return videos, nil
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

func (repo *VideoRepository) GetPublic(ctx context.Context, id uuid.UUID) (model.Video, error) {
	const op = "postgres.VideoRepository.GetPublic"

	query := `
		SELECT videos.*, authors.*
		FROM videos
		JOIN users AS authors ON videos.author_id = authors.id
		WHERE videos.id = $1 AND videos.is_public = true
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

func (repo *VideoRepository) Update(ctx context.Context, id uuid.UUID, dto repository.UpdateVideoDTO) error {
	const op = "postgres.VideoRepository.Update"

	var (
		counter int   = 1
		args    []any = []any{id}
		query   strings.Builder
	)
	_, _ = query.WriteString("UPDATE videos SET ")

	if dto.Title != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("title = $%d, ", counter))
		args = append(args, *dto.Title)
	}

	if dto.Description != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("description = $%d, ", counter))
		args = append(args, *dto.Description)
	}

	if dto.ThumbnailPath != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("thumbnail_path = $%d, ", counter))
		args = append(args, *dto.ThumbnailPath)
	}

	if dto.VideoPath != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("video_path = $%d, ", counter))
		args = append(args, *dto.VideoPath)
	}

	if dto.Public != nil {
		counter++
		_, _ = query.WriteString(fmt.Sprintf("is_public = $%d, ", counter))
		args = append(args, *dto.Public)
	}

	_, _ = query.WriteString("updated_at = now() WHERE id = $1")

	_, err := repo.db.ExecContext(ctx, query.String(), args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (repo *VideoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "postgres.VideoRepository.Delete"

	query := `DELETE FROM videos WHERE id = $1`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
