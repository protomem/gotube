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

func (acc *CommentAccessor) AllByVideoID(ctx context.Context, videoID entity.ID, opts port.FindOptions) ([]entity.Comment, error) {
	comments, err := acc.commentDao.SelectByVideoID(ctx, videoID, database.SelectOptions(opts))
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Comment{}, nil
		}

		return nil, err
	}

	authorIDs := make([]entity.ID, 0, len(comments))
	for _, comment := range comments {
		authorIDs = append(authorIDs, comment.AuthorID)
	}

	authors, err := acc.userDao.SelectByIDs(ctx, authorIDs)
	if err != nil {
		if database.IsNoRows(err) {
			return []entity.Comment{}, nil
		}

		return nil, err
	}

	return mapCommentEntriesAndUserEntriesToCommentEntities(comments, authors), nil
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

func mapCommentEntriesAndUserEntriesToCommentEntities(comments []database.CommentEntry, users []database.UserEntry) []entity.Comment {
	entities := make([]entity.Comment, 0, len(comments))
	for _, comment := range comments {
		for _, user := range users {
			if comment.AuthorID == user.ID {
				entities = append(entities, mapCommentEntryAndUserEntryToCommentEntity(comment, user))
				break
			}
		}
	}
	return entities
}
