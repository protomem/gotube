package database

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID uuid.UUID `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	FromUserID uuid.UUID `db:"from_user_id" json:"-"`
	FromUser   User      `db:"-" json:"fromUser"`

	ToUserID uuid.UUID `db:"to_user_id" json:"-"`
	ToUser   User      `db:"-" json:"toUser"`
}
