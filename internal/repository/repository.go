package repository

import "database/sql"

type Repository struct {
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{}
}
