package database

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/domain/model"
)

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
	if row.Err() != nil {
		if IsNoRows(row.Err()) {
			return model.Video{}, model.ErrVideoNotFound
		}

		return model.Video{}, row.Err()
	}

	if err := db.videoScan(row, &video); err != nil {
		return model.Video{}, err
	}

	return video, nil
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
