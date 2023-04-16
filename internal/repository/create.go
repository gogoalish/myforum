package repository

import "database/sql"

func Create(DB *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL,
		name TEXT NOT NULL,
		password TEXT NOT NULL
		);`
	_, err := DB.Exec(query)
	return err
}
