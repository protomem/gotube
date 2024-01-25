package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type VideoEntry struct {
	ID uuid.UUID `db:"id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Title       string `db:"title"`
	Description string `db:"description"`

	ThumbnailPath string `db:"thumbnail_path"`
	VideoPath     string `db:"video_path"`

	AuthorID uuid.UUID `db:"author_id"`

	Public bool   `db:"is_public"`
	Views  uint64 `db:"views"`
}

type VideoDAO struct {
	db *DB
}

func (db *DB) VideoDAO() *VideoDAO {
	return &VideoDAO{
		db: db,
	}
}

func (dao *VideoDAO) GetByID(ctx context.Context, id uuid.UUID) (VideoEntry, error) {
	const op = "database.VideoDAO.GetByID"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT * FROM videos WHERE id = $1 LIMIT 1`
	args := []any{id}

	var video VideoEntry

	if err := dao.db.QueryRowxContext(ctx, query, args...).StructScan(&video); err != nil {
		return VideoEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

type InsertVideoDTO struct {
	Title         string
	Description   string
	ThumbnailPath string
	VideoPath     string
	AuthorID      uuid.UUID
	Public        bool
}

func (dao *VideoDAO) Insert(ctx context.Context, dto InsertVideoDTO) (uuid.UUID, error) {
	const op = "database.VideoDAO.Insert"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO videos(title, description, thumbnail_path, video_path, author_id, is_public) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	args := []any{dto.Title, dto.Description, dto.ThumbnailPath, dto.VideoPath, dto.AuthorID, dto.Public}

	var id uuid.UUID

	if err := dao.db.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
