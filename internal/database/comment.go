package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var ErrCommentNotFound = NewModelError(ErrNotFound, "comment")

type Comment struct {
	ID uuid.UUID `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	Body string `db:"body" json:"body"`

	VideoID  uuid.UUID `db:"video_id" json:"videoId"`
	AuthorID uuid.UUID `db:"author_id" json:"-"`

	Author User `db:"author" json:"author"`
}

func (db *DB) FindCommentsByVideoID(ctx context.Context, videoID uuid.UUID, opts FindOptions) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
		SELECT comments.*, authors.* FROM comments
		JOIN users AS authors ON comments.author_id = authors.id
		WHERE comments.video_id = $3
		LIMIT $1 OFFSET $2
	`
	args := []any{opts.Limit, opts.Offset, videoID}

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []Comment{}, nil
		}

		return []Comment{}, err
	}
	defer func() { _ = rows.Close() }()

	comments := make([]Comment, 0, opts.Limit)

	for rows.Next() {
		var comment Comment
		if err := db.commmentScan(rows, &comment); err != nil {
			return []Comment{}, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (db *DB) GetComment(ctx context.Context, id uuid.UUID) (Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
		SELECT comments.*, authors.* FROM comments
		JOIN users AS authors ON comments.author_id = authors.id
		WHERE comments.id = $1
		LIMIT 1
	`

	var comment Comment

	row := db.QueryRowxContext(ctx, query, id)
	if err := db.commmentScan(row, &comment); err != nil {
		if IsNoRows(err) {
			return Comment{}, ErrCommentNotFound
		}

		return Comment{}, err
	}

	return comment, nil
}

type InsertCommentDTO struct {
	Body     string
	VideoID  uuid.UUID
	AuthorID uuid.UUID
}

func (db *DB) InsertComment(ctx context.Context, dto InsertCommentDTO) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `INSERT INTO comments (body, video_id, author_id) VALUES ($1, $2, $3) RETURNING id`
	args := []any{dto.Body, dto.VideoID, dto.AuthorID}

	var id uuid.UUID

	if err := db.
		QueryRowxContext(ctx, query, args...).
		Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (db *DB) DeleteComment(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `DELETE FROM comments WHERE id = $1`
	args := []any{id}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (db *DB) commmentScan(s Scanner, comment *Comment) error {
	return s.Scan(
		&comment.ID,
		&comment.CreatedAt, &comment.UpdatedAt,
		&comment.Body,
		&comment.VideoID, &comment.AuthorID,

		&comment.Author.ID,
		&comment.Author.CreatedAt, &comment.Author.UpdatedAt,
		&comment.Author.Nickname, &comment.Author.Password,
		&comment.Author.Email,
		&comment.Author.AvatarPath, &comment.Author.Description,
	)
}
