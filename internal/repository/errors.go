package repository

import (
	"database/sql"
	"errors"
	"strings"
)

var (
	ErrNoRows               = errors.New("no rows")
	ErrForeignKeyConstraint = errors.New("foreign key constraint failed")
	ErrUserExists           = errors.New("user already exists")
)

func isNoRowsError(err error) bool {
	if err == nil {
		return false
	}
	return err == sql.ErrNoRows
}

func isAlreadyExists(err error) bool {
	if err == nil {
		return false
	}

	return strings.HasPrefix(err.Error(), "UNIQUE constraint failed:")
}

func isForeignKeyConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "FOREIGN KEY constraint failed")
}
