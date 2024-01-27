package database

import (
	"time"

	"github.com/google/uuid"
)

type CommentEntry struct {
	ID uuid.UUID `db:"id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Content string `db:"content"`

	AuthorID uuid.UUID `db:"author_id"`
	VideoID  uuid.UUID `db:"video_id"`
}

type CommentDAO struct {
	db *DB
}

func (db *DB) CommentDAO() *CommentDAO {
	return &CommentDAO{
		db: db,
	}
}
