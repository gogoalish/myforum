package repository

import (
	"database/sql"

	"forum/internal/models"
)

type PostRepo struct {
	DB *sql.DB
}

type Posts interface {
	Create(UserID int, title, content string) error
	All() ([]*models.Post, error)
}

func (r *PostRepo) Create(userID int, title, content string) error {
	query := `INSERT INTO posts
	VALUES(NULL, ?, ?, ?)`
	if _, err := r.DB.Exec(query, userID, title, content); err != nil {
		return nil
	}
	return nil
}

func (r *PostRepo) All() ([]*models.Post, error) {
	posts := []*models.Post{}
	query := `SELECT * FROM posts`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := &models.Post{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, err
}
