package service

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

var _ Comment = (*CommentImpl)(nil)

type (
	CreateCommentDTO struct {
		Message  string
		VideoID  model.ID
		AuthorID model.ID
	}
)

type (
	Comment interface {
		FindByVideo(ctx context.Context, videoID model.ID, opts FindOptions) ([]model.Comment, error)
		Create(ctx context.Context, dto CreateCommentDTO) (model.Comment, error)
		Delete(ctx context.Context, id model.ID) error
	}

	CommentImpl struct {
		repo repository.Comment
	}
)

func NewComment(repo repository.Comment) *CommentImpl {
	return &CommentImpl{
		repo: repo,
	}
}

func (s *CommentImpl) FindByVideo(ctx context.Context, videoID model.ID, opts FindOptions) ([]model.Comment, error) {
	const op = "service.Comment.FindByVideo"

	comments, err := s.repo.FindByVideo(ctx, videoID, repository.FindOptions(opts))
	if err != nil {
		return []model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}

func (s *CommentImpl) Create(ctx context.Context, dto CreateCommentDTO) (model.Comment, error) {
	const op = "service.Comment.Create"

	// TODO: Add validation

	id, err := s.repo.Create(ctx, repository.CreateCommentDTO(dto))
	if err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	comment, err := s.repo.Get(ctx, id)
	if err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (s *CommentImpl) Delete(ctx context.Context, id model.ID) error {
	const op = "service.Comment.Delete"

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
