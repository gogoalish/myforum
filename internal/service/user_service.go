package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type Users interface {
	SignUp(models.User) error
	Get(email, name string) (models.User, error)
}

type UserService struct {
	repository.Users
}

func NewUserService(r repository.Users) Users {
	return &UserService{r}
}

func (u *UserService) SignUp(m models.User) error {
	PasswordCrypt(&m)
	u.Users.SignUp(m)
	return nil
}

func (u *UserService) Get(email, name string) (models.User, error) {
	m, err := u.Users.Get(email, name)
	return m, err
}
