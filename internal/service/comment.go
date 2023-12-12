package service

import "github.com/protomem/gotube/internal/repository"

type (
	Comment interface{}

	CommentImpl struct {
		repo     repository.Comment
		userServ User
	}
)

func NewComment(repo repository.Comment, userServ User) *CommentImpl {
	return &CommentImpl{
		repo:     repo,
		userServ: userServ,
	}
}
