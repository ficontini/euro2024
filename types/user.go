package types

import "golang.org/x/crypto/bcrypt"

type User struct {
	Email             string
	EncryptedPassword string
}

func NewUser(email, password string) (User, error) {
	var user User
	encpwd, err := generateEncryptedPassword(password)
	if err != nil {
		return user, err
	}
	{
		user.Email = email
		user.EncryptedPassword = encpwd
	}
	return user, nil
}

func generateEncryptedPassword(password string) (string, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encpw), nil
}
