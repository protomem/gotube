package sqlite

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsKeyConflict(err error) bool {
	var sqlErr sqlite3.Error
	return errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique
}
