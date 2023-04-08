package repository

import "database/sql"

type Repository struct {
	Users Users
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Users: NewUserRepo(db),
	}
}
