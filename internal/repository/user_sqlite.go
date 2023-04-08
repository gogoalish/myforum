package repository

import (
	"database/sql"

	"forum/internal/models"
)

type UserRepo struct {
	DB *sql.DB
}

type Users interface {
	Create(models.User) error
	Get(email, name string) (models.User, error)
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

func (u *UserRepo) Create(m models.User) error {
	query := `INSERT INTO users (email, name, password)
	VALUES(?, ?, ?);`
	if _, err := u.DB.Exec(query, m.Email, m.Name, m.Password); err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) Get(email, name string) (models.User, error) {
	query := `SELECT * FROM users
	WHERE ? = email AND ? = name`
	var user models.User
	err := u.DB.QueryRow(query, email, name).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	return user, err
}
