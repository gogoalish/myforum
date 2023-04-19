package repository

import "database/sql"

type Repo struct {
	Users
}

func New(db *sql.DB) *Repo {
	return &Repo{
		Users: NewUserRepo(db),
	}
}
