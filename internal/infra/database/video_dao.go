package database

import (
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
