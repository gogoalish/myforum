package service

import (
	"forum/internal/models"
	"forum/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Users interface {
	SignUp(models.User) error
	SignIn(email, password string) (models.User, error)
}

type UserService struct {
	repo repository.Users
}

func NewUserService(r *repository.Repo) Users {
	return &UserService{r.Users}
}

func (u *UserService) SignUp(m models.User) error {
	Encrypt(&m)
	u.repo.SignUp(m)
	return nil
}

func (u *UserService) SignIn(email, password string) (models.User, error) {
	m, err := u.repo.SignIn(email)
	if err != nil {
		return m, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	return m, err
}
