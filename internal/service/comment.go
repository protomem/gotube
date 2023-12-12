package service

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ Comment = (*CommentImpl)(nil)

type CreateCommentDTO struct {
	Message  string
	AuthorID model.ID
	VideoID  model.ID
}

func (dto CreateCommentDTO) Validate() error {
	return nil
}

type (
	Comment interface {
		Create(ctx context.Context, dto CreateCommentDTO) (model.Comment, error)
	}

	CommentImpl struct {
		repo      repository.Comment
		userServ  User
		videoServ Video
	}
)

func NewComment(repo repository.Comment, userServ User, videoServ Video) *CommentImpl {
	return &CommentImpl{
		repo:      repo,
		userServ:  userServ,
		videoServ: videoServ,
	}
}

func (s *CommentImpl) Create(ctx context.Context, dto CreateCommentDTO) (model.Comment, error) {
	const op = "service:Comment.Create"

	if _, err := s.videoServ.Get(ctx, dto.VideoID); err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	commentID, err := s.repo.Create(ctx, repository.CreateCommentDTO(dto))
	if err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	comment, err := s.repo.Get(ctx, commentID)
	if err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	comment.Author, err = s.userServ.Get(ctx, comment.Author.ID)
	if err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}
