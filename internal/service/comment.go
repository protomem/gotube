package service

import "github.com/protomem/gotube/internal/repository"

var _ Comment = (*CommentImpl)(nil)

type (
	Comment interface{}

	CommentImpl struct {
		repo repository.Comment
	}
)

func NewComment(repo repository.Comment) *CommentImpl {
	return &CommentImpl{
		repo: repo,
	}
}
