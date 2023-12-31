package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type ModelError struct {
	InternalError error
	Model         string
}

func NewModelError(err error, model string) *ModelError {
	return &ModelError{
		InternalError: err,
		Model:         model,
	}
}

func (e *ModelError) Error() string {
	return fmt.Sprintf("%s: %s", e.Model, e.InternalError)
}

func (e *ModelError) Is(err error) bool {
	return errors.Is(err, e.InternalError)
}

func (e *ModelError) Unwrap() error {
	return e.InternalError
}

func (e *ModelError) As(target any) bool {
	if _, ok := target.(*ModelError); ok {
		return true
	}
	if errors.As(e.InternalError, target) {
		return true
	}
	return false
}

func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsKeyConflict(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}
