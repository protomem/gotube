package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/protomem/gotube/internal/module/user/model"
)

var ErrSubscriptionAlreadyExists = errors.New("subscription already exists")

type Subscription struct {
	ID uuid.UUID `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	FromUser usermodel.User `json:"fromUser"`
	ToUser   usermodel.User `json:"toUser"`
}
