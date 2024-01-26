package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type VideoEntry struct {
	ID uuid.UUID `db:"id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Title       string `db:"title"`
	Description string `db:"description"`

	ThumbnailPath string `db:"thumbnail_path"`
	VideoPath     string `db:"video_path"`

	AuthorID uuid.UUID `db:"author_id"`

	Public bool   `db:"is_public"`
	Views  uint64 `db:"views"`
}

type VideoDAO struct {
	db *DB
}

func (db *DB) VideoDAO() *VideoDAO {
	return &VideoDAO{
		db: db,
	}
}

func (dao *VideoDAO) GetByID(ctx context.Context, id uuid.UUID) (VideoEntry, error) {
	const op = "database.VideoDAO.GetByID"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT * FROM videos WHERE id = $1 LIMIT 1`
	args := []any{id}

	var video VideoEntry

	if err := dao.db.QueryRowxContext(ctx, query, args...).StructScan(&video); err != nil {
		return VideoEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

type InsertVideoDTO struct {
	Title         string
	Description   string
	ThumbnailPath string
	VideoPath     string
	AuthorID      uuid.UUID
	Public        bool
}

func (dao *VideoDAO) Insert(ctx context.Context, dto InsertVideoDTO) (uuid.UUID, error) {
	const op = "database.VideoDAO.Insert"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO videos(title, description, thumbnail_path, video_path, author_id, is_public) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	args := []any{dto.Title, dto.Description, dto.ThumbnailPath, dto.VideoPath, dto.AuthorID, dto.Public}

	var id uuid.UUID

	if err := dao.db.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

type UpdateVideoDTO struct {
	Title         *string
	Description   *string
	ThumbnailPath *string
	VideoPath     *string
	Public        *bool
}

func (dao *VideoDAO) Update(ctx context.Context, id uuid.UUID, dto UpdateVideoDTO) error {
	const op = "database.VideoDAO.Update"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	counter := 1
	query := `UPDATE videos SET updated_at = now()`
	args := []any{id}

	if dto.Title != nil {
		counter++
		query += `, title = $` + strconv.Itoa(counter)
		args = append(args, *dto.Title)
	}
	if dto.Description != nil {
		counter++
		query += `, description = $` + strconv.Itoa(counter)
		args = append(args, *dto.Description)
	}
	if dto.ThumbnailPath != nil {
		counter++
		query += `, thumbnail_path = $` + strconv.Itoa(counter)
		args = append(args, *dto.ThumbnailPath)
	}
	if dto.VideoPath != nil {
		counter++
		query += `, video_path = $` + strconv.Itoa(counter)
		args = append(args, *dto.VideoPath)
	}
	if dto.Public != nil {
		counter++
		query += `, is_public = $` + strconv.Itoa(counter)
		args = append(args, *dto.Public)
	}

	query += ` WHERE id = $1`

	if _, err := dao.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (dao *VideoDAO) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "database.VideoDAO.Delete"

	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `DELETE FROM videos WHERE id = $1`
	args := []any{id}

	if _, err := dao.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
