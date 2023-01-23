package repository

import "database/sql"

type Repository struct {
	User User
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: NewUser(db),
	}
}
