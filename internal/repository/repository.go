package repository

import "database/sql"

type Repos struct {
	Users
}

func New(db *sql.DB) *Repos {
	return &Repos{
		Users: NewUserRepo(db),
	}
}
