package models

type SignForm struct {
	Email     string
	Name      string
	Password  string
	Validator map[string]string
}
