package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/protomem/gotube/internal/database/postgres"
	"github.com/protomem/gotube/internal/module/video/model"
	"github.com/protomem/gotube/internal/module/video/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.VideoRepository = (*VideoRepository)(nil)

type VideoRepository struct {
	logger  logging.Logger
	db      *postgres.DB
	builder squirrel.StatementBuilderType
}

func NewVideoRepository(logger logging.Logger, db *postgres.DB) *VideoRepository {
	return &VideoRepository{
		logger:  logger.With("repository", "video", "repositoryType", "postgres"),
		db:      db,
		builder: postgres.Builder(),
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

func (r *VideoRepository) FindAllVideosByUserIDAndSortByNew(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.Video, error) {
	const op = "VideoRepository.FindAllVideosByUserIDAndSortByNew"
	var err error

	query, args, err := r.builder.
		Select("videos.*, users.*").
		From("videos").
		Where(squirrel.Eq{"videos.user_id": userID}).
		Join("users ON users.id = videos.user_id").
		OrderBy("videos.created_at DESC").
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

func (r *VideoRepository) FindAllVideosWherePublicAndSortByNew(
	ctx context.Context,
	options repository.FindAllVideosOptions,
) ([]model.Video, uint64, error) {
	const op = "VideoRepository.FindAllPublicNewVideos"
	var err error

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: begin: %w", op, err)
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				r.logger.Error("failed to cancel transaction", "error", err, "operation", op)
			}
		}
	}()

	selectQuery, selectArgs, err := r.builder.
		Select("videos.*, users.*").
		From("videos").
		Where(squirrel.Eq{"videos.is_public": true}).
		Join("users ON users.id = videos.user_id").
		OrderBy("videos.created_at DESC").
		Limit(options.Limit).
		Offset(options.Offset).
		ToSql()
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := tx.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.Video{}, 0, nil
		}

		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
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
			return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
		}

		videos = append(videos, video)
	}

	if len(videos) == 0 {
		return []model.Video{}, 0, nil
	}

	countQuery, countArgs, err := r.builder.
		Select("COUNT(*)").
		From("videos").
		Where(squirrel.Eq{"videos.is_public": true}).
		ToSql()
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	var count uint64
	err = tx.QueryRow(ctx, countQuery, countArgs...).Scan(&count)
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: commit: %w", op, err)
	}

	return videos, count, nil
}

func (r *VideoRepository) FindAllVideosWherePublicAndSortByPopular(
	ctx context.Context,
	options repository.FindAllVideosOptions,
) ([]model.Video, uint64, error) {
	const op = "VideoRepository.FindAllPublicPopularVideos"
	var err error

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: begin: %w", op, err)
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				r.logger.Error("failed to cancel transaction", "error", err, "operation", op)
			}
		}
	}()

	selectQuery, selectArgs, err := r.builder.
		Select("videos.*, users.*").
		From("videos").
		Where(squirrel.Eq{"videos.is_public": true}).
		Join("users ON users.id = videos.user_id").
		OrderBy("videos.views DESC").
		Limit(options.Limit).
		Offset(options.Offset).
		ToSql()
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := tx.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.Video{}, 0, nil
		}

		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
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
			return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
		}

		videos = append(videos, video)
	}

	if len(videos) == 0 {
		return []model.Video{}, 0, nil
	}

	countQuery, countArgs, err := r.builder.
		Select("COUNT(*)").
		From("videos").
		Where(squirrel.Eq{"videos.is_public": true}).
		ToSql()
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	var count uint64
	err = tx.QueryRow(ctx, countQuery, countArgs...).Scan(&count)
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.Video{}, 0, fmt.Errorf("%s: commit: %w", op, err)
	}

	return videos, count, nil
}

func (r *VideoRepository) FindAllVideosByUserIDsAndWherePublicAndSortByNew(
	ctx context.Context,
	userIDs []uuid.UUID,
	options repository.FindAllVideosOptions,
) ([]model.Video, error) {
	const op = "VideoRepository.FindAllPublicNewVideosByUserIDs"
	var err error

	query, args, err := r.builder.
		Select("videos.*, users.*").
		From("videos").
		Where(squirrel.Eq{"videos.is_public": true}).
		Where("videos.user_id = ANY(?::UUID[])", userIDs).
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

func (r *VideoRepository) FindAllVideosLikeByTitleAndWherePublicAndSortByNew(
	ctx context.Context,
	title string,
	options repository.FindAllVideosOptions,
) ([]model.Video, error) {
	const op = "VideoRepository.FindAllPublicNewVideosLikeByTitle"
	var err error

	query, args, err := r.builder.
		Select("videos.*, users.*").
		From("videos").
		Where(squirrel.Eq{"videos.is_public": true}).
		Where("videos.title LIKE ?", "%"+title+"%").
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
		if postgres.IsDuplicateKeyError(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrVideoAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *VideoRepository) IncrementVideoView(ctx context.Context, id uuid.UUID) error {
	const op = "VideoRepository.IncrementVideoView"
	var err error

	query := `
        UPDATE videos SET views = views + 1 WHERE id = $1
    `

	args := []any{
		id,
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
