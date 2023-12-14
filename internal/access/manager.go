package access

import "context"

const (
	Read   Action = "read"
	Edit   Action = "edit"
	Delete Action = "delete"
)

type Action string

const (
	User         Object = "user"
	Subscription Object = "subscription"
	Video        Object = "video"
	Rating       Object = "rating"
	Comment      Object = "comment"
)

type Object string

type Attribute struct {
	Key    string
	Values string
}

type Manager interface {
	Close(ctx context.Context) error
}
