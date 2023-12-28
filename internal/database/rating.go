package database

import (
	"time"

	"github.com/google/uuid"
)

type Rating struct {
	ID uuid.UUID `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	VideoID uuid.UUID `db:"video_id" json:"videoId"`
	UserID  uuid.UUID `db:"user_id" json:"userId"`

	Liked bool `db:"is_liked" json:"isLiked"`
}
