package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
)

type (
	CreateCommentDTO struct {
		Content  string
		AuthorID uuid.UUID
		VideoID  uuid.UUID
	}

	Comment interface {
		FindByVideoID(ctx context.Context, videoID uuid.UUID) ([]model.Comment, error)

		Create(ctx context.Context, dto CreateCommentDTO) (model.Comment, error)

		Delete(ctx context.Context, id uuid.UUID) error
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

func (serv *CommentImpl) FindByVideoID(ctx context.Context, videoID uuid.UUID) ([]model.Comment, error) {
	const op = "service.Comment.FindByVideoID"

	comments, err := serv.repo.FindByVideoID(ctx, videoID)
	if err != nil {
		return []model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}

func (serv *CommentImpl) Create(ctx context.Context, dto CreateCommentDTO) (model.Comment, error) {
	const op = "service.Comment.Create"
	var err error

	// TODO: Valiate ...

	id, err := serv.repo.Create(ctx, repository.CreateCommentDTO(dto))
	if err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	comment, err := serv.repo.Get(ctx, id)
	if err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (serv *CommentImpl) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "service.Comment.Delete"

	err := serv.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
