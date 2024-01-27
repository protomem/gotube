package handler

type Comment struct {
	*Base
}

func NewComment() *Comment {
	return &Comment{NewBase()}
}
