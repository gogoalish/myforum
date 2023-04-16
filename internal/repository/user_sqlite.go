package repository

import (
	"database/sql"

	"forum/internal/models"
)

type Users interface {
	SignUp(models.User) error
	Get(email, name string) (models.User, error)
}

type UserRepo struct {
	*sql.DB
}

func NewUserRepo(db *sql.DB) Users {
	return &UserRepo{db}
}

func (u *UserRepo) SignUp(m models.User) error {
	query := `INSERT INTO users (email, name, password)
	VALUES(?, ?, ?);`
	if _, err := u.Exec(query, m.Email, m.Name, m.Password); err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) Get(email, name string) (models.User, error) {
	query := `SELECT * FROM users
	WHERE ? = email AND ? = name`
	var user models.User
	err := u.QueryRow(query, email, name).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	return user, err
}
