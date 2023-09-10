package model

import (
	"time"

	"github.com/google/uuid"
	usermodel "github.com/protomem/gotube/internal/module/user/model"
)

type Video struct {
	ID uuid.UUID `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Title       string `json:"title"`
	Description string `json:"description"`

	ThumbnailPath string `json:"thumbnailPath"`
	VideoPath     string `json:"videoPath"`

	Public bool `json:"isPublic"`

	Views int `json:"views"`

	User usermodel.User `json:"user"`
}
