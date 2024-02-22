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

var _ repository.Video = (*Video)(nil)

type videoEntry struct {
	ID            string
	CreatedAt     int64
	UpdatedAt     int64
	Title         string
	Description   string
	ThumbnailPath string
	VideoPath     string
	AuthorID      string
	Author        userEntry
	Public        bool
	Views         int64
}

type Video struct {
	logger logging.Logger
	db     database.DB
}

func NewVideo(logger logging.Logger, db database.DB) *Video {
	return &Video{
		logger: logger.With("repository", "sqlite/video"),
		db:     db,
	}
}

func (r *Video) FindSortByCreatedAtWherePublic(ctx context.Context, opts repository.FindOptions) ([]model.Video, error) {
	const op = "repository.Video.FindSortByCreatedAtWherePublic"

	query := `
		SELECT videos.*, authors.* FROM videos
		JOIN users AS authors ON videos.author_id = authors.id
		WHERE videos.is_public > 0
		ORDER BY videos.created_at DESC
		LIMIT ? OFFSET ?
	`
	args := []any{opts.Limit, opts.Offset}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		if sqlite.IsNoRows(err) {
			return []model.Video{}, nil
		}

		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = rows.Close() }()

	videos := make([]model.Video, 0, opts.Limit)
	for rows.Next() {
		video, err := r.scan(rows)
		if err != nil {
			return []model.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		videos = append(videos, video)
	}

	return videos, nil
}

func (r *Video) FindSortByViewsWherePublic(ctx context.Context, opts repository.FindOptions) ([]model.Video, error) {
	const op = "repository.Video.FindSortByCreatedAtWherePublic"

	query := `
		SELECT videos.*, authors.* FROM videos
		JOIN users AS authors ON videos.author_id = authors.id
		WHERE videos.is_public > 0
		ORDER BY videos.views DESC
		LIMIT ? OFFSET ?
	`
	args := []any{opts.Limit, opts.Offset}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		if sqlite.IsNoRows(err) {
			return []model.Video{}, nil
		}

		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = rows.Close() }()

	videos := make([]model.Video, 0, opts.Limit)
	for rows.Next() {
		video, err := r.scan(rows)
		if err != nil {
			return []model.Video{}, fmt.Errorf("%s: %w", op, err)
		}

		videos = append(videos, video)
	}

	return videos, nil
}

func (r *Video) Get(ctx context.Context, id model.ID) (model.Video, error) {
	const op = "repository.Video.Get"

	query := `
		SELECT videos.*, authors.* FROM videos
		JOIN users AS authors ON videos.author_id = authors.id
		WHERE videos.id = ?
		LIMIT 1
	`
	args := []any{id.String()}

	row := r.db.QueryRow(ctx, query, args...)
	video, err := r.scan(row)
	if err != nil {
		if sqlite.IsNoRows(err) {
			return model.Video{}, fmt.Errorf("%s: %w", op, model.ErrVideoNotFound)
		}

		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (r *Video) Create(ctx context.Context, dto repository.CreateVideoDTO) (model.ID, error) {
	const op = "repository.Video.Create"

	id, err := uuid.NewRandom()
	if err != nil {
		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}
	now := time.Now()

	query := `
		INSERT INTO videos (id, created_at, updated_at, title, description, thumbnail_path, video_path, author_id, is_public)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	args := []any{id.String(), now.Unix(), now.Unix(), dto.Title, dto.Description, dto.ThumbnailPath, dto.VideoPath, dto.AuthorID.String(), dto.Public}

	if err := r.db.Exec(ctx, query, args...); err != nil {
		if sqlite.IsKeyConflict(err) {
			return model.ID{}, fmt.Errorf("%s: %w", op, model.ErrVideoExists)
		}

		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Video) Update(ctx context.Context, id model.ID, dto repository.UpdateVideoDTO) error {
	const op = "repository.Video.Update"

	now := time.Now()

	query := `UPDATE videos SET updated_at = ?`
	args := []any{now.Unix()}

	if dto.Title != nil {
		query += `, title = ?`
		args = append(args, *dto.Title)
	}
	if dto.Description != nil {
		query += `, description = ?`
		args = append(args, *dto.Description)
	}
	if dto.ThumbnailPath != nil {
		query += `, thumbnail_path = ?`
		args = append(args, *dto.ThumbnailPath)
	}
	if dto.VideoPath != nil {
		query += `, video_path = ?`
		args = append(args, *dto.VideoPath)
	}
	if dto.Public != nil {
		query += `, is_public = ?`
		args = append(args, *dto.Public)
	}

	query += ` WHERE id = ?`
	args = append(args, id.String())

	if err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Video) Delete(ctx context.Context, id model.ID) error {
	const op = "repository.Video.Delete"

	query := `DELETE FROM videos WHERE id = ?`
	args := []any{id.String()}

	if err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Video) scan(s database.Scanner) (model.Video, error) {
	var entry videoEntry
	if err := s.Scan(
		&entry.ID, &entry.CreatedAt, &entry.UpdatedAt,
		&entry.Title, &entry.Description,
		&entry.ThumbnailPath, &entry.VideoPath,
		&entry.AuthorID, &entry.Public,
		&entry.Views,

		&entry.Author.ID, &entry.Author.CreatedAt, &entry.Author.UpdatedAt,
		&entry.Author.Nickname, &entry.Author.Password,
		&entry.Author.Email, &entry.Author.Verified,
		&entry.Author.AvatarPath, &entry.Author.Description,
	); err != nil {
		return model.Video{}, err
	}

	id, err := uuid.Parse(entry.ID)
	if err != nil {
		return model.Video{}, err
	}

	createdAt := time.Unix(entry.CreatedAt, 0)
	updatedAt := time.Unix(entry.UpdatedAt, 0)

	authorID, err := uuid.Parse(entry.Author.ID)
	if err != nil {
		return model.Video{}, err
	}

	authorCreatedAt := time.Unix(entry.Author.CreatedAt, 0)
	authorUpdatedAt := time.Unix(entry.Author.UpdatedAt, 0)

	return model.Video{
		Model: model.Model{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Title:         entry.Title,
		Description:   entry.Description,
		ThumbnailPath: entry.ThumbnailPath,
		VideoPath:     entry.VideoPath,
		Public:        entry.Public,
		Views:         entry.Views,
		Author: model.User{
			Model: model.Model{
				ID:        authorID,
				CreatedAt: authorCreatedAt,
				UpdatedAt: authorUpdatedAt,
			},
			Nickname:    entry.Author.Nickname,
			Password:    entry.Author.Password,
			Email:       entry.Author.Email,
			Verified:    entry.Author.Verified,
			AvatarPath:  entry.Author.AvatarPath,
			Description: entry.Author.Description,
		},
	}, nil
}
