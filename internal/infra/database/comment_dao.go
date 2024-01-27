package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CommentEntry struct {
	ID uuid.UUID `db:"id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Content string `db:"content"`

	AuthorID uuid.UUID `db:"author_id"`
	VideoID  uuid.UUID `db:"video_id"`
}

type CommentDAO struct {
	db *DB
}

func (db *DB) CommentDAO() *CommentDAO {
	return &CommentDAO{
		db: db,
	}
}

func (dao *CommentDAO) SelectByVideoID(ctx context.Context, videoID uuid.UUID, opts SelectOptions) ([]CommentEntry, error) {
	const op = "database.CommentDAO.SelectByVideoID"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT * FROM comments WHERE video_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	args := []any{videoID, opts.Limit, opts.Offset}

	comments := make([]CommentEntry, 0, opts.Limit)

	if err := dao.db.SelectContext(ctx, &comments, query, args...); err != nil {
		return []CommentEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}

func (dao *CommentDAO) GetByID(ctx context.Context, id uuid.UUID) (CommentEntry, error) {
	const op = "database.CommentDAO.GetByID"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT * FROM comments WHERE id = $1 LIMIT 1`
	args := []any{id}

	var comment CommentEntry

	if err := dao.db.QueryRowxContext(ctx, query, args...).StructScan(&comment); err != nil {
		return CommentEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

type InsertCommentDTO struct {
	Content  string
	AuthorID uuid.UUID
	VideoID  uuid.UUID
}

func (dao *CommentDAO) Insert(ctx context.Context, dto InsertCommentDTO) (uuid.UUID, error) {
	const op = "database.CommentDAO.Insert"

	query := `INSERT INTO comments(content, author_id, video_id) VALUES ($1, $2, $3) RETURNING id`
	args := []any{dto.Content, dto.AuthorID, dto.VideoID}

	var id uuid.UUID

	if err := dao.db.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
