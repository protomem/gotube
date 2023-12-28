package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var ErrRatingAlreadyExists = NewModelError(ErrAlreadyExists, "rating")

type Rating struct {
	ID uuid.UUID `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	VideoID uuid.UUID `db:"video_id" json:"videoId"`
	UserID  uuid.UUID `db:"user_id" json:"userId"`

	Liked bool `db:"is_liked" json:"isLiked"`
}

func (db *DB) CountRatingsWithLikedByVideoID(ctx context.Context, videoID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT COUNT(*) FROM ratings WHERE video_id = $1 AND is_liked = true`
	args := []any{videoID}

	var count int

	if err := db.
		QueryRowxContext(ctx, query, args...).
		Scan(&count); err != nil {
		if IsNoRows(err) {
			return 0, nil
		}

		return 0, err
	}

	return count, nil
}

func (db *DB) CountRatingsWithNotLikedByVideoID(ctx context.Context, videoID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT COUNT(*) FROM ratings WHERE video_id = $1 AND is_liked = false`
	args := []any{videoID}

	var count int

	if err := db.
		QueryRowxContext(ctx, query, args...).
		Scan(&count); err != nil {
		if IsNoRows(err) {
			return 0, nil
		}

		return 0, err
	}

	return count, nil
}

func (db *DB) GetRatingByVideoIDAndUserID(ctx context.Context, videoID, userID uuid.UUID) (Rating, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `SELECT * FROM ratings WHERE video_id = $1 AND user_id = $2 LIMIT 1`
	args := []any{videoID, userID}

	var rating Rating

	if err := db.
		QueryRowxContext(ctx, query, args...).
		StructScan(&rating); err != nil {
		if IsNoRows(err) {
			return Rating{}, nil
		}

		return Rating{}, err
	}

	return rating, nil
}

type InsertRatingDTO struct {
	VideoID uuid.UUID
	UserID  uuid.UUID
	Liked   bool
}

func (db *DB) InsertRating(ctx context.Context, dto InsertRatingDTO) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `INSERT INTO ratings(video_id, user_id, is_liked) VALUES ($1, $2, $3) RETURNING id`
	args := []any{dto.VideoID, dto.UserID, dto.Liked}

	var id uuid.UUID

	if err := db.
		QueryRowxContext(ctx, query, args...).
		Scan(&id); err != nil {
		if IsKeyConflict(err) {
			return uuid.Nil, ErrRatingAlreadyExists
		}

		return uuid.Nil, err
	}

	return id, nil
}

type UpdateRatingDTO struct {
	Liked bool
}

func (db *DB) UpdateRating(ctx context.Context, id uuid.UUID, dto UpdateRatingDTO) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `UPDATE ratings SET updated_at = now(), is_liked = $1 WHERE id = $2`
	args := []any{dto.Liked, id}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteRating(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	query := `DELETE FROM ratings WHERE id = $1`
	args := []any{id}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
