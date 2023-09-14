package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/module/video/model"
	"github.com/protomem/gotube/internal/module/video/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.VideoRepository = (*VideoRepository)(nil)

type VideoRepository struct {
	logger  logging.Logger
	db      *database.DB
	builder squirrel.StatementBuilderType
}

func NewVideoRepository(logger logging.Logger, db *database.DB) *VideoRepository {
	return &VideoRepository{
		logger:  logger.With("repository", "video", "repositoryType", "postgres"),
		db:      db,
		builder: database.Builder(),
	}
}

func (r *VideoRepository) FindOneVideo(ctx context.Context, videoID uuid.UUID) (model.Video, error) {
	const op = "VideoRepository.FindOneVideo"
	var err error

	query, args, err := r.builder.
		Select("videos.*, users.*").
		From("videos").
		Where(squirrel.Eq{"videos.id": videoID}).
		Join("users ON users.id = videos.user_id").
		Limit(1).
		ToSql()
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	var video model.Video
	err = r.db.
		QueryRow(ctx, query, args...).
		Scan(
			&video.ID, &video.CreatedAt, &video.UpdatedAt,
			&video.Title, &video.Description,
			&video.ThumbnailPath, &video.VideoPath,
			&video.Public, &video.Views,
			&video.User.ID,

			&video.User.ID, &video.User.CreatedAt, &video.User.UpdatedAt,
			&video.User.Nickname, &video.User.Password,
			&video.User.Email, &video.User.Verified,
			&video.User.AvatarPath, &video.User.Description,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Video{}, fmt.Errorf("%s: %w", op, model.ErrVideoNotFound)
		}

		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (r *VideoRepository) FindAllPublicNewVideos(
	ctx context.Context,
	options repository.FindAllVideosOptions,
) ([]model.Video, error) {
	const op = "VideoRepository.FindAllPublicNewVideos"
	var err error

	query, args, err := r.builder.
		Select("videos.*, users.*").
		From("videos").
		Where(squirrel.Eq{"videos.is_public": true}).
		Join("users ON users.id = videos.user_id").
		OrderBy("videos.created_at DESC").
		Limit(options.Limit).
		Offset(options.Offset).
		ToSql()
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.Video{}, nil
		}

		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var videos []model.Video
	for rows.Next() {
		var video model.Video
		err = rows.Scan(
			&video.ID, &video.CreatedAt, &video.UpdatedAt,
			&video.Title, &video.Description,
			&video.ThumbnailPath, &video.VideoPath,
			&video.Public, &video.Views,
			&video.User.ID,

			&video.User.ID, &video.User.CreatedAt, &video.User.UpdatedAt,
			&video.User.Nickname, &video.User.Password,
			&video.User.Email, &video.User.Verified,
			&video.User.AvatarPath, &video.User.Description,
		)
		if err != nil {
			return []model.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		videos = append(videos, video)
	}

	if len(videos) == 0 {
		return []model.Video{}, nil
	}

	return videos, nil
}

func (r *VideoRepository) CreateVideo(ctx context.Context, dto repository.CreateVideoDTO) (uuid.UUID, error) {
	const op = "VideoRepository.CreateVideo"
	var err error

	query, args, err := r.builder.
		Insert("videos").
		Columns("title", "description", "thumbnail_path", "video_path", "user_id").
		Values(dto.Title, dto.Description, dto.ThumbnailPath, dto.VideoPath, dto.UserID).
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
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrVideoAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
