package repository

import "database/sql"

type Repo struct {
	Users
	Posts
}

func New(db *sql.DB) *Repo {
	return &Repo{
		Users: NewUserRepo(db),
		Posts: &PostRepo{db},
	}
}
