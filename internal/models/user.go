package models

type User struct {
	ID       int
	Email    string
	Name     string
	Password string
	Token    *string
	Expires  *string
}
