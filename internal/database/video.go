package database

import (
	"time"

	"github.com/google/uuid"
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
	Author   User      `db:"-" json:"author"`
}
