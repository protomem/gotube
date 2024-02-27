package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/gotube/internal/database"
	"github.com/protomem/gotube/internal/database/sqlite"
	"github.com/protomem/gotube/internal/model"
	"github.com/protomem/gotube/internal/repository"
	"github.com/protomem/gotube/pkg/logging"
)

var _ repository.Comment = (*Comment)(nil)

type commentEntry struct {
	ID        string
	CreatedAt int64
	UpdatedAt int64
	Message   string
	VideoID   string
	AuthorID  string
	Author    userEntry
}

type Comment struct {
	logger logging.Logger
	db     database.DB
}

func NewComment(logger logging.Logger, db database.DB) *Comment {
	return &Comment{
		logger: logger.With("repository", "sqlite/comment"),
		db:     db,
	}
}

func (r *Comment) Get(ctx context.Context, id model.ID) (model.Comment, error) {
	const op = "repository.Comment.Find"

	query := `
		SELECT comments.*, authors.* FROM comments
		JOIN users AS authors ON comments.author_id = authors.id
		WHERE comments.id = ?
		LIMIT 1
	`
	args := []any{id.String()}

	row := r.db.QueryRow(ctx, query, args...)
	comment, err := r.scan(row)
	if err != nil {
		if sqlite.IsNoRows(err) {
			return model.Comment{}, fmt.Errorf("%s: %w", op, model.ErrCommentNotFound)
		}

		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (r *Comment) Create(ctx context.Context, dto repository.CreateCommentDTO) (model.ID, error) {
	const op = "repository.Comment.Create"

	id, err := uuid.NewRandom()
	if err != nil {
		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}
	now := time.Now()

	query := `INSERT INTO comments (id, created_at, updated_at, message, video_id, author_id) VALUES (?, ?, ?, ?, ?, ?)`
	args := []any{id.String(), now.Unix(), now.Unix(), dto.Message, dto.VideoID.String(), dto.AuthorID.String()}

	if err := r.db.Exec(ctx, query, args...); err != nil {
		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Comment) scan(s database.Scanner) (model.Comment, error) {
	var entry commentEntry
	if err := s.Scan(
		&entry.ID, &entry.CreatedAt, &entry.UpdatedAt,
		&entry.Message, &entry.VideoID, &entry.AuthorID,

		&entry.Author.ID, &entry.Author.CreatedAt, &entry.Author.UpdatedAt,
		&entry.Author.Nickname, &entry.Author.Password,
		&entry.Author.Email, &entry.Author.Verified,
		&entry.Author.AvatarPath, &entry.Author.Description,
	); err != nil {
		return model.Comment{}, err
	}

	id, err := uuid.Parse(entry.ID)
	if err != nil {
		return model.Comment{}, err
	}

	createdAt := time.Unix(entry.CreatedAt, 0)
	updatedAt := time.Unix(entry.UpdatedAt, 0)

	videoID, err := uuid.Parse(entry.VideoID)
	if err != nil {
		return model.Comment{}, err
	}

	authorID, err := uuid.Parse(entry.Author.ID)
	if err != nil {
		return model.Comment{}, err
	}

	authorCreatedAt := time.Unix(entry.Author.CreatedAt, 0)
	authorUpdatedAt := time.Unix(entry.Author.UpdatedAt, 0)

	return model.Comment{
		Model: model.Model{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Message: entry.Message,
		VideoID: videoID,
		Author: model.User{
			Model: model.Model{
				ID:        authorID,
				CreatedAt: authorCreatedAt,
				UpdatedAt: authorUpdatedAt,
			},
			Nickname:    entry.Author.Nickname,
			Password:    entry.Author.Password,
			Email:       entry.Author.Email,
			Verified:    entry.Author.Verified,
			AvatarPath:  entry.Author.AvatarPath,
			Description: entry.Author.Description,
		},
	}, nil
}
