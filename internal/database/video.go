package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrVideoNotFound      = NewModelError(ErrNotFound, "video")
	ErrVideoAlreadyExists = NewModelError(ErrAlreadyExists, "video")
)

type Video struct {
	ID uuid.UUID `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`

	ThumbnailPath string `db:"thumbnail_path" json:"thumbnailPath"`
	VideoPath     string `db:"video_path" json:"videoPath"`

	Public bool   `db:"is_public" json:"isPublic"`
	Views  uint64 `db:"views" json:"views"`

	AuthorID uuid.UUID `db:"author_id" json:"-"`
	Author   User      `db:"author" json:"author"`
}

func (db *DB) GetVideo(ctx context.Context, id uuid.UUID) (Video, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        SELECT videos.*, author.* FROM videos 
        JOIN users AS author ON videos.author_id = author.id
        WHERE videos.id = $1 LIMIT 1
    `
	args := []any{id}

	var video Video

	row := db.QueryRowxContext(ctx, query, args...)
	if err := db.videoScan(row, &video); err != nil {
		if IsNoRows(err) {
			return Video{}, ErrVideoNotFound
		}

		return Video{}, err
	}

	return video, nil
}

type InsertVideoDTO struct {
	Title         string
	Description   string
	ThumbnailPath string
	VideoPath     string
	Public        bool
	AuthorID      uuid.UUID
}

func (db *DB) InsertVideo(ctx context.Context, dto InsertVideoDTO) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        INSERT INTO videos (title, description, thumbnail_path, video_path, is_public, author_id) 
        VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING id
    `
	args := []any{dto.Title, dto.Description, dto.ThumbnailPath, dto.VideoPath, dto.Public, dto.AuthorID}

	var id uuid.UUID

	if err := db.
		QueryRowxContext(ctx, query, args...).
		Scan(&id); err != nil {
		if IsKeyConflict(err) {
			return uuid.Nil, ErrVideoAlreadyExists
		}

		return uuid.Nil, err
	}

	return id, nil
}

func (db *DB) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := "DELETE FROM videos WHERE id = $1"
	args := []any{id}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (db *DB) videoScan(row *sqlx.Row, video *Video) error {
	return row.Scan(
		&video.ID,
		&video.CreatedAt, &video.UpdatedAt,
		&video.Title, &video.Description,
		&video.ThumbnailPath, &video.VideoPath,
		&video.Public, &video.Views,
		&video.AuthorID,

		&video.Author.ID,
		&video.Author.CreatedAt, &video.Author.UpdatedAt,
		&video.Author.Nickname, &video.Author.Password,
		&video.Author.Email,
		&video.Author.AvatarPath, &video.Author.Description,
	)
}
