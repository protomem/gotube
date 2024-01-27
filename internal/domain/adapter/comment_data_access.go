package adapter

import (
	"context"

	"github.com/protomem/gotube/internal/domain/entity"
	"github.com/protomem/gotube/internal/domain/port"
	"github.com/protomem/gotube/internal/infra/database"
)

var _ port.CommentAccessor = (*CommentAccessor)(nil)

type CommentAccessor struct {
	commentDao *database.CommentDAO
	userDao    *database.UserDAO
}

func NewCommentAccessor(db *database.DB) *CommentAccessor {
	return &CommentAccessor{
		commentDao: db.CommentDAO(),
		userDao:    db.UserDAO(),
	}
}

func (acc *CommentAccessor) ByID(ctx context.Context, id entity.ID) (entity.Comment, error) {
	comment, err := acc.commentDao.GetByID(ctx, id)
	if err != nil {
		if database.IsNoRows(err) {
			return entity.Comment{}, entity.ErrCommentNotFound
		}

		return entity.Comment{}, err
	}

	author, err := acc.userDao.GetByID(ctx, comment.AuthorID)
	if err != nil {
		if database.IsNoRows(err) {
			return entity.Comment{}, entity.ErrCommentNotFound
		}

		return entity.Comment{}, err
	}

	return mapCommentEntryAndUserEntryToCommentEntity(comment, author), nil
}

var _ port.CommentMutator = (*CommentMutator)(nil)

type CommentMutator struct {
	dao *database.CommentDAO
}

func NewCommentMutator(db *database.DB) *CommentMutator {
	return &CommentMutator{db.CommentDAO()}
}

func (mut *CommentMutator) Insert(ctx context.Context, dto port.InsertCommentDTO) (entity.ID, error) {
	id, err := mut.dao.Insert(ctx, database.InsertCommentDTO(dto))
	if err != nil {
		return entity.ID{}, err
	}

	return id, nil
}

func mapCommentEntryAndUserEntryToCommentEntity(comment database.CommentEntry, user database.UserEntry) entity.Comment {
	return entity.Comment{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Content:   comment.Content,
		Author:    entity.User(user),
		VideoID:   comment.VideoID,
	}
}
