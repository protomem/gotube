package database

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var (
	ErrVideoNotFound      = NewModelError(ErrNotFound, "video")
	ErrVideoAlreadyExists = NewModelError(ErrAlreadyExists, "video")
)

type Video struct {
	ID uuid.UUID `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`

	ThumbnailPath string `db:"thumbnail_path" json:"thumbnailPath"`
	VideoPath     string `db:"video_path" json:"videoPath"`

	Public bool   `db:"is_public" json:"isPublic"`
	Views  uint64 `db:"views" json:"views"`

	AuthorID uuid.UUID `db:"author_id" json:"-"`
	Author   User      `db:"author" json:"author"`
}

func (db *DB) FindPublicVideosSortByCreatedAt(ctx context.Context, opts FindOptions) ([]Video, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        SELECT videos.*, authors.* FROM videos 
        JOIN users AS authors ON videos.author_id = authors.id
        WHERE videos.is_public = true
        ORDER BY videos.created_at DESC
        LIMIT $1 OFFSET $2
    `
	args := []any{opts.Limit, opts.Offset}

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []Video{}, nil
		}

		return []Video{}, err
	}
	defer func() { _ = rows.Close() }()

	videos := make([]Video, 0, opts.Limit)

	for rows.Next() {
		var video Video
		if err = db.videoScan(rows, &video); err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *DB) FindPublicVideosSortByViews(ctx context.Context, opts FindOptions) ([]Video, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        SELECT videos.*, authors.* FROM videos
        JOIN users AS authors ON videos.author_id = authors.id 
        WHERE videos.is_public = true
        ORDER BY videos.views DESC
        LIMIT $1 OFFSET $2
    `
	args := []any{opts.Limit, opts.Offset}

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []Video{}, nil
		}

		return []Video{}, err
	}
	defer func() { _ = rows.Close() }()

	videos := make([]Video, 0, opts.Limit)

	for rows.Next() {
		var video Video
		if err = db.videoScan(rows, &video); err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *DB) FindPublicVideosByAuthorNicknameSortByCreatedAt(
	ctx context.Context,
	authorNickname string,
	opts FindOptions,
) ([]Video, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        SELECT videos.*, authors.* FROM videos
        JOIN users AS authors ON videos.author_id = authors.id
        WHERE authors.nickname = $3 AND videos.is_public = true
        ORDER BY videos.created_at DESC
        LIMIT $1 OFFSET $2
    `
	args := []any{opts.Limit, opts.Offset, authorNickname}

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []Video{}, nil
		}

		return []Video{}, err
	}
	defer func() { _ = rows.Close() }()

	videos := make([]Video, 0, opts.Limit)

	for rows.Next() {
		var video Video
		if err = db.videoScan(rows, &video); err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *DB) FindVideosByAuthorNicknameSortByCreatedAt(
	ctx context.Context,
	authorNickname string,
	opts FindOptions,
) ([]Video, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        SELECT videos.*, authors.* FROM videos
        JOIN users AS authors ON videos.author_id = authors.id
        WHERE authors.nickname = $3
        ORDER BY videos.created_at DESC
        LIMIT $1 OFFSET $2
    `
	args := []any{opts.Limit, opts.Offset, authorNickname}

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []Video{}, nil
		}

		return []Video{}, err
	}
	defer func() { _ = rows.Close() }()

	videos := make([]Video, 0, opts.Limit)

	for rows.Next() {
		var video Video
		if err = db.videoScan(rows, &video); err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *DB) FindPublicVideosLikeByTitle(ctx context.Context, likeTitle string, opts FindOptions) ([]Video, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        SELECT videos.*, authors.* FROM videos
        JOIN users AS authors ON videos.author_id = authors.id
        WHERE lower(videos.title) LIKE '%' || lower($3) || '%' AND videos.is_public = true
        ORDER BY videos.created_at DESC
        LIMIT $1 OFFSET $2
    `
	args := []any{opts.Limit, opts.Offset, likeTitle}

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []Video{}, nil
		}

		return []Video{}, err
	}
	defer func() { _ = rows.Close() }()

	videos := make([]Video, 0, opts.Limit)

	for rows.Next() {
		var video Video
		if err = db.videoScan(rows, &video); err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *DB) FindPublicVideosLikeByTitleAndAuthorNickname(
	ctx context.Context,
	likeTitle, authorNickname string,
	opts FindOptions,
) ([]Video, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        SELECT videos.*, authors.* FROM videos
        JOIN users AS authors ON videos.author_id = authors.id
        WHERE lower(videos.title) LIKE '%' || lower($3) || '%' AND authors.nickname = $4 AND videos.is_public = true
        ORDER BY videos.created_at DESC
        LIMIT $1 OFFSET $2
    `
	args := []any{opts.Limit, opts.Offset, likeTitle, authorNickname}

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []Video{}, nil
		}

		return []Video{}, err
	}
	defer func() { _ = rows.Close() }()

	videos := make([]Video, 0, opts.Limit)

	for rows.Next() {
		var video Video
		if err = db.videoScan(rows, &video); err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *DB) FindVideosLikeByTitleAndAuthorNickname(
	ctx context.Context,
	likeTitle, authorNickname string,
	opts FindOptions,
) ([]Video, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        SELECT videos.*, authors.* FROM videos
        JOIN users AS authors ON videos.author_id = authors.id
        WHERE lower(videos.title) LIKE '%' || lower($3) || '%' AND authors.nickname = $4
        ORDER BY videos.created_at DESC
        LIMIT $1 OFFSET $2
    `
	args := []any{opts.Limit, opts.Offset, likeTitle, authorNickname}

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []Video{}, nil
		}

		return []Video{}, err
	}
	defer func() { _ = rows.Close() }()

	videos := make([]Video, 0, opts.Limit)

	for rows.Next() {
		var video Video
		if err = db.videoScan(rows, &video); err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *DB) FindPublicVideosByAuthorIDsSortByCreatedAt(ctx context.Context, authorIDs []uuid.UUID, opts FindOptions) ([]Video, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
		SELECT videos.*, authors.* FROM videos
		JOIN users AS authors ON videos.author_id = authors.id
		WHERE authors.id = ANY($3::uuid[]) AND videos.is_public = true
		ORDER BY videos.created_at DESC
		LIMIT $1 OFFSET $2
	`
	args := []any{opts.Limit, opts.Offset, authorIDs}

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		if IsNoRows(err) {
			return []Video{}, nil
		}

		return []Video{}, err
	}
	defer func() { _ = rows.Close() }()

	videos := make([]Video, 0, opts.Limit)
	for rows.Next() {
		var video Video
		if err = db.videoScan(rows, &video); err != nil {
			return []Video{}, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (db *DB) GetVideo(ctx context.Context, id uuid.UUID) (Video, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        SELECT videos.*, author.* FROM videos 
        JOIN users AS author ON videos.author_id = author.id
        WHERE videos.id = $1 LIMIT 1
    `
	args := []any{id}

	var video Video

	row := db.QueryRowxContext(ctx, query, args...)
	if err := db.videoScan(row, &video); err != nil {
		if IsNoRows(err) {
			return Video{}, ErrVideoNotFound
		}

		return Video{}, err
	}

	return video, nil
}

type InsertVideoDTO struct {
	Title         string
	Description   string
	ThumbnailPath string
	VideoPath     string
	Public        bool
	AuthorID      uuid.UUID
}

func (db *DB) InsertVideo(ctx context.Context, dto InsertVideoDTO) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `
        INSERT INTO videos (title, description, thumbnail_path, video_path, is_public, author_id) 
        VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING id
    `
	args := []any{dto.Title, dto.Description, dto.ThumbnailPath, dto.VideoPath, dto.Public, dto.AuthorID}

	var id uuid.UUID

	if err := db.
		QueryRowxContext(ctx, query, args...).
		Scan(&id); err != nil {
		if IsKeyConflict(err) {
			return uuid.Nil, ErrVideoAlreadyExists
		}

		return uuid.Nil, err
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

func (db *DB) UpdateVideo(ctx context.Context, id uuid.UUID, dto UpdateVideoDTO) error {
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

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := "DELETE FROM videos WHERE id = $1"
	args := []any{id}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (db *DB) videoScan(s Scanner, video *Video) error {
	return s.Scan(
		&video.ID,
		&video.CreatedAt, &video.UpdatedAt,
		&video.Title, &video.Description,
		&video.ThumbnailPath, &video.VideoPath,
		&video.Public, &video.Views,
		&video.AuthorID,

		&video.Author.ID,
		&video.Author.CreatedAt, &video.Author.UpdatedAt,
		&video.Author.Nickname, &video.Author.Password,
		&video.Author.Email,
		&video.Author.AvatarPath, &video.Author.Description,
	)
}
