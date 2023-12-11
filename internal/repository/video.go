package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/pkg/logging"
)

var _ Video = (*VideoImpl)(nil)

type CreateVideoDTO struct {
	Title         string
	Description   string
	ThumbnailPath string
	VideoPath     string
	AuthorID      model.ID
	Public        bool
}

type (
	Video interface {
		Find(ctx context.Context, opts FindOptions) ([]model.Video, error)
		FindOrderByViews(ctx context.Context, opts FindOptions) ([]model.Video, error)
		Get(ctx context.Context, id model.ID) (model.Video, error)
		Create(ctx context.Context, dto CreateVideoDTO) (model.ID, error)
	}

	VideoImpl struct {
		logger logging.Logger
		pdb    *pgxpool.Pool
	}
)

func NewVideo(logger logging.Logger, pdb *pgxpool.Pool) *VideoImpl {
	return &VideoImpl{
		logger: logger.With("repository", "video", "repositoryType", "postgres"),
		pdb:    pdb,
	}
}

func (r *VideoImpl) Find(ctx context.Context, opts FindOptions) ([]model.Video, error) {
	const op = "repository:Video.Find"

	query := `
        SELECT videos.*, authors.*
        FROM videos 
            JOIN users AS authors 
            ON videos.author_id = authors.id 
        WHERE videos.is_public = true
        ORDER BY videos.created_at DESC
        LIMIT $1
        OFFSET $2
    `
	args := []any{opts.Limit, opts.Offset}

	rows, err := r.pdb.Query(ctx, query, args...)
	if err != nil {
		if IsPgNotFound(err) {
			return []model.Video{}, nil
		}

		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() { rows.Close() }()

	videos := make([]model.Video, 0, opts.Limit)
	for rows.Next() {
		var video model.Video
		if err := r.scan(rows, &video); err != nil {
			return []model.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		videos = append(videos, video)
	}

	if len(videos) == 0 {
		return []model.Video{}, nil
	}

	return videos, nil
}

func (r *VideoImpl) FindOrderByViews(ctx context.Context, opts FindOptions) ([]model.Video, error) {
	const op = "repository:Video.FindOrderByViews"

	query := `
        SELECT videos.*, authors.*
        FROM videos 
            JOIN users AS authors 
            ON videos.author_id = authors.id 
        WHERE videos.is_public = true
        ORDER BY videos.views DESC
        LIMIT $1
        OFFSET $2
    `
	args := []any{opts.Limit, opts.Offset}

	rows, err := r.pdb.Query(ctx, query, args...)
	if err != nil {
		if IsPgNotFound(err) {
			return []model.Video{}, nil
		}

		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() { rows.Close() }()

	videos := make([]model.Video, 0, opts.Limit)
	for rows.Next() {
		var video model.Video
		if err := r.scan(rows, &video); err != nil {
			return []model.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		videos = append(videos, video)
	}

	if len(videos) == 0 {
		return []model.Video{}, nil
	}

	return videos, nil
}

func (r *VideoImpl) Get(ctx context.Context, id model.ID) (model.Video, error) {
	const op = "repository:Video.Get"

	query := `
        SELECT videos.*, authors.* 
        FROM videos 
            JOIN users AS authors 
            ON videos.author_id = authors.id 
        WHERE videos.id = $1 
        LIMIT 1
    `
	args := []any{id}

	row := r.pdb.QueryRow(ctx, query, args...)

	var video model.Video
	if err := r.scan(row, &video); err != nil {
		if IsPgNotFound(err) {
			return model.Video{}, fmt.Errorf("%s: %w", op, model.ErrVideoNotFound)
		}

		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (r *VideoImpl) Create(ctx context.Context, dto CreateVideoDTO) (model.ID, error) {
	const op = "repository:Video.Create"

	query := `
        INSERT INTO videos (title, description, thumbnail_path, video_path, author_id, is_public) 
        VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING id
    `
	args := []any{
		dto.Title, dto.Description,
		dto.ThumbnailPath, dto.VideoPath,
		dto.AuthorID, dto.Public,
	}

	row := r.pdb.QueryRow(ctx, query, args...)

	var id model.ID
	if err := row.Scan(&id); err != nil {
		if IsPgDuplicateKey(err) {
			return model.ID{}, fmt.Errorf("%s: %w", op, model.ErrVideoAlreadyExists)
		}

		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *VideoImpl) scan(row Scanable, video *model.Video) error {
	return row.Scan(
		&video.ID,
		&video.CreatedAt, &video.UpdatedAt,
		&video.Title, &video.Description,
		&video.ThumbnailPath, &video.VideoPath,
		&video.Author.ID,
		&video.Public, &video.Views,

		&video.Author.ID,
		&video.Author.CreatedAt, &video.Author.UpdatedAt,
		&video.Author.Nickname, &video.Author.Password,
		&video.Author.Email, &video.Author.Verified,
		&video.Author.AvatarPath, &video.Author.Description,
	)
}
