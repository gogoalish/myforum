package service

import (
	"forum/internal/models"
	"forum/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Users
}

type Users interface {
	Create(models.User) error
	Get(email, name string) (models.User, error)
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		NewUserService(r.Users),
	}
}

type UserService struct {
	repository.Users
}

func NewUserService(r repository.Users) Users {
	return &UserService{r}
}

func (u *UserService) Create(m models.User) error {
	PasswordCrypt(&m)
	u.Users.Create(m)
	return nil
}

func (u *UserService) Get(email, name string) (models.User, error) {
	m, err := u.Users.Get(email, name)
	return m, err
}

func PasswordCrypt(m *models.User) {
	password := []byte(m.Password)
	crypted, _ := bcrypt.GenerateFromPassword(password, 3)
	m.Password = string(crypted)
}
