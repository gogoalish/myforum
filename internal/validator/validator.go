package validator

import (
	"net/mail"
	"regexp"
	"unicode/utf8"

	"forum/internal/models"
)

const (
	MsgEmailExists        = "Email already exists"
	MsgNameExists         = "Name already exists"
	MsgInvalidEmail       = "Write correct email"
	MsgInvalidName        = "Write correct name"
	MsgUserNotFound       = "User not found"
	MsgNotCorrectPassword = "Password is not correct"
)

func GetErrMsgs(m models.User) map[string]string {
	errmsgs := make(map[string]string)
	if !isValidEmail(m.Email) {
		errmsgs["email"] = MsgInvalidEmail
	}
	if !isValidName(m.Name) {
		errmsgs["name"] = MsgInvalidName
	}
	return errmsgs
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidName(name string) bool {
	length := utf8.RuneCountInString(name)
	if length < 5 || length > 15 {
		return false
	}
	usernameConvention := "^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$"
	re, _ := regexp.Compile(usernameConvention)
	return re.MatchString(name)
}
