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
		FindByVideoID(ctx context.Context, videoID model.ID) ([]model.Comment, error)
		Create(ctx context.Context, dto CreateCommentDTO) (model.Comment, error)
		Delete(ctx context.Context, id model.ID) error
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

func (s *CommentImpl) FindByVideoID(ctx context.Context, videoID model.ID) ([]model.Comment, error) {
	const op = "service:Comment.FindByVideoID"

	videos, err := s.repo.FindByVideoID(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userIDs := make([]model.ID, 0, len(videos))
	for _, video := range videos {
		userIDs = append(userIDs, video.Author.ID)
	}

	users, err := s.userServ.FindByIDs(ctx, userIDs...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for i, video := range videos {
		for _, user := range users {
			if video.Author.ID == user.ID {
				videos[i].Author = user
			}
		}
	}

	return videos, nil
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

func (s *CommentImpl) Delete(ctx context.Context, id model.ID) error {
	const op = "service:Comment.Delete"

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
