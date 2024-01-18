package database

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/domain/model"
)

func (db *DB) FindVideosSortByCreatedAt(ctx context.Context, opts FindOptions) ([]model.Video, error) {
	const op = "database.FindVideosSortByCreatedAt"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	videos, err := db.findVideosByFieldWithSortBy(ctx, Field{Name: "is_public", Value: true}, Field{Name: "created_at", Value: SortByDesc}, opts)
	if err != nil {
		return []model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return videos, nil
}

func (db *DB) GetVideo(ctx context.Context, id model.ID) (model.Video, error) {
	const op = "database.GetVideo"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	video, err := db.getVideoByField(ctx, Field{Name: "id", Value: id})
	if err != nil {
		return model.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

type InsertVideoDTO struct {
	Title         string
	Description   string
	ThumbnailPath string
	VideoPath     string
	Public        bool
	AuthorID      model.ID
}

func (db *DB) InsertVideo(ctx context.Context, dto InsertVideoDTO) (model.ID, error) {
	const op = "database.InsertVideo"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO videos(title, description, thumbnail_path, video_path, is_public, author_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	args := []any{dto.Title, dto.Description, dto.ThumbnailPath, dto.VideoPath, dto.Public, dto.AuthorID}

	var id model.ID

	if err := db.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		if IsKeyConflict(err) {
			return model.ID{}, fmt.Errorf("%s: %w", op, model.ErrVideoAlreadyExists)
		}

		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

type UpdateVideoDTO struct {
	Title         *string
	Description   *string
	ThumbnailPath *string
	VideoPath     *string
	Public        *bool
}

func (db *DB) UpdateVideo(ctx context.Context, id model.ID, dto UpdateVideoDTO) error {
	const op = "database.UpdateVideo"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	fields := make([]Field, 0, 5)

	if dto.Title != nil {
		fields = append(fields, Field{Name: "title", Value: *dto.Title})
	}
	if dto.Description != nil {
		fields = append(fields, Field{Name: "description", Value: *dto.Description})
	}
	if dto.ThumbnailPath != nil {
		fields = append(fields, Field{Name: "thumbnail_path", Value: *dto.ThumbnailPath})
	}
	if dto.VideoPath != nil {
		fields = append(fields, Field{Name: "video_path", Value: *dto.VideoPath})
	}
	if dto.Public != nil {
		fields = append(fields, Field{Name: "is_public", Value: *dto.Public})
	}

	if err := db.updateVideoByField(ctx, Field{Name: "id", Value: id}, fields); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (db *DB) DeleteVideo(ctx context.Context, id model.ID) error {
	const op = "database.DeleteVideo"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	if err := db.deleteVideoByField(ctx, Field{Name: "id", Value: id}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (db *DB) findVideosByFieldWithSortBy(ctx context.Context, byField Field, sortBy Field, opts FindOptions) ([]model.Video, error) {
	baseQuery := `
        SELECT videos.*, authors.* FROM videos
        JOIN users AS authors ON videos.author_id = authors.id
        WHERE videos.%s = $1
        ORDER BY videos.%s %s
        LIMIT $2 OFFSET $3
    `
	query := fmt.Sprintf(baseQuery, byField.Name, sortBy.Name, sortBy.Value)
	args := []any{byField.Value, opts.Limit, opts.Offset}

	videos := make([]model.Video, 0)

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []model.Video{}, nil
		}

		return []model.Video{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var video model.Video
		if err := db.videoScan(rows, &video); err != nil {
			return []model.Video{}, err
		}

		videos = append(videos, video)
	}

	return videos, nil
}

func (db *DB) getVideoByField(ctx context.Context, field Field) (model.Video, error) {
	baseQuery := `
		SELECT videos.*, authors.* FROM videos 
		JOIN users AS authors ON videos.author_id = authors.id 
		WHERE videos.%s = $1 LIMIT 1
	`
	query := fmt.Sprintf(baseQuery, field.Name)
	args := []any{field.Value}

	var video model.Video

	row := db.QueryRowxContext(ctx, query, args...)
	if err := db.videoScan(row, &video); err != nil {
		if IsNoRows(err) {
			return model.Video{}, model.ErrVideoNotFound
		}

		return model.Video{}, err
	}

	return video, nil
}

func (db *DB) updateVideoByField(ctx context.Context, byFiled Field, fields []Field) error {
	counter := 1
	query := `UPDATE videos SET updated_at = now()`
	args := []any{byFiled.Value}

	for _, f := range fields {
		counter++
		query += fmt.Sprintf(`, %s = $%d`, f.Name, counter)
		args = append(args, f.Value)
	}

	query += fmt.Sprintf(` WHERE %s = $%d`, byFiled.Name, 1)

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (db *DB) deleteVideoByField(ctx context.Context, field Field) error {
	query := fmt.Sprintf(`DELETE FROM videos WHERE %s = $1`, field.Name)
	args := []any{field.Value}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (db *DB) videoScan(s Scanner, video *model.Video) error {
	return s.Scan(
		&video.ID,
		&video.CreatedAt, &video.UpdatedAt,
		&video.Title, &video.Description,
		&video.ThumbnailPath, &video.VideoPath,
		&video.AuthorID,
		&video.Views, &video.Public,

		&video.Author.ID,
		&video.Author.CreatedAt, &video.Author.UpdatedAt,
		&video.Author.Nickname, &video.Author.Password,
		&video.Author.Email,
		&video.Author.AvatarPath, &video.Author.Description,
	)
}
