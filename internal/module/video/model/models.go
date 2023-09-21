package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/protomem/gotube/internal/module/user/model"
)

var (
	ErrVideoNotFound      = errors.New("video not found")
	ErrVideoAlreadyExists = errors.New("video already exists")
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

const (
	Like    RatingType = "like"
	Dislike RatingType = "dislike"
)

type RatingType string

var ErrRatingAlreadyExists = errors.New("rating already exists")

type Rating struct {
	ID uuid.UUID `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Type RatingType `json:"type"`

	VideoID uuid.UUID `json:"videoID"`
	UserID  uuid.UUID `json:"userID"`
}
