package repository

import (
	"database/sql"
	"errors"

	"forum/internal/models"
)

type Users interface {
	SignUp(models.User) error
	SignIn(email string) (models.User, error)
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

func (u *UserRepo) SignIn(email string) (models.User, error) {
	query := `SELECT * FROM users
	WHERE ? = email`
	var user models.User
	err := u.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Token)
	if errors.Is(err, sql.ErrNoRows) {
		return user, models.ErrNoRecord
	}
	return user, err
}
