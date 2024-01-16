package database

import (
	"context"
	"fmt"

	"github.com/protomem/gotube/internal/domain/model"
)

func (db *DB) GetComment(ctx context.Context, id model.ID) (model.Comment, error) {
	const op = "database.GetComment"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	comment, err := db.getCommentByField(ctx, Field{Name: "id", Value: id})
	if err != nil {
		return model.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

type InsertCommentDTO struct {
	Content  string
	AuthorID model.ID
	VideoID  model.ID
}

func (db *DB) InsertComment(ctx context.Context, dto InsertCommentDTO) (model.ID, error) {
	const op = "database.InsertComment"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        INSERT INTO comments (content, author_id, video_id) 
        VALUES ($1, $2, $3) 
        RETURNING id
    `
	args := []any{dto.Content, dto.AuthorID, dto.VideoID}

	var id model.ID

	if err := db.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		return model.ID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (db *DB) getCommentByField(ctx context.Context, field Field) (model.Comment, error) {
	baseQuery := `
		SELECT comments.*, authors.* FROM comments
		JOIN users AS authors ON comments.author_id = authors.id 
		WHERE comments.%s = $1 LIMIT 1
	`
	query := fmt.Sprintf(baseQuery, field.Name)
	args := []any{field.Value}

	var comment model.Comment

	row := db.QueryRowxContext(ctx, query, args...)
	if err := db.commentScan(row, &comment); err != nil {
		if IsNoRows(err) {
			return model.Comment{}, model.ErrCommentNotFound
		}

		return model.Comment{}, err
	}

	return comment, nil
}

func (db *DB) commentScan(s Scanner, comment *model.Comment) error {
	return s.Scan(
		&comment.ID,
		&comment.CreatedAt, &comment.UpdatedAt,
		&comment.Content,
		&comment.AuthorID, &comment.VideoID,

		&comment.Author.ID,
		&comment.Author.CreatedAt, &comment.Author.UpdatedAt,
		&comment.Author.Nickname, &comment.Author.Password,
		&comment.Author.Email,
		&comment.Author.AvatarPath, &comment.Author.Description,
	)
}
