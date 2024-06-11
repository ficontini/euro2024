package types

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	minNameLen     = 3
	minPasswordLen = 7
	firstName      = "fistName"
	lastName       = "lastName"
	password       = "password"
	email          = "email"
)

type User struct {
	FirstName         string
	LastName          string
	Email             string
	EncryptedPassword string
}

func NewUserFromParams(params UserParams) (*User, error) {
	encpwd, err := generateEncryptedPassword(params.Password)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: encpwd,
	}, nil
}

func generateEncryptedPassword(password string) (string, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encpw), nil
}

type UserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p UserParams) Validate() map[string]string {
	errors := make(map[string]string)
	if len(p.FirstName) < minNameLen {
		errors[firstName] = fmt.Sprintf("%s length must be at least %d characters", firstName, minNameLen)
	}
	if len(p.LastName) < minNameLen {
		errors[lastName] = fmt.Sprintf("%s length must be at least %d characters", lastName, minNameLen)
	}
	if !isPasswordValid(p.Password) {
		errors[password] = fmt.Sprintf("%s (%s) is invalid", password, p.Password)
	}
	if !isEmailValid(p.Email) {
		errors[email] = fmt.Sprintf("%s (%s) is invalid", email, p.Email)
	}
	return errors
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

// TODO:
func isPasswordValid(password string) bool {
	return len(password) >= minPasswordLen
}
