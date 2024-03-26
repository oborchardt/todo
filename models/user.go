package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system.
type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	password string
}

// SetPassword sets the password for the user.
func (user *User) SetPassword(password string) {
	user.password = password
}

// GetPasswordHash generates and returns the bcrypt hash of the user's password.
func (user *User) GetPasswordHash() ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.password), 14)
	return bytes, err
}

// CheckPassword compares the provided hashed password with the user's actual password.
// It returns nil if the passwords match, otherwise, it returns an error indicating the mismatch.
func (user *User) CheckPassword(hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(user.password))
}

// UserLogin represents the login credentials for a user.
type UserLogin struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
}

// UserFromLogin constructs a User object from the provided UserLogin credentials.
// It validates the presence of both username and password fields in the UserLogin object.
// If either field is missing, it returns an error indicating the missing field.
func UserFromLogin(login UserLogin) (User, error) {
	var user User
	if login.Name == nil || login.Password == nil {
		return user, errors.New("User name or password are missing.")
	}
	user.Name = *login.Name
	user.SetPassword(*login.Password)
	return user, nil
}
