package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	password string
}

func (user *User) SetPassword(password string) {
	user.password = password
}

func (user *User) GetPasswordHash() ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.password), 14)
	return bytes, err
}

func (user *User) CheckPassword(hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(user.password))
}

type UserLogin struct {
	UserId   *int    `json:"userId"`
	Password *string `json:"password"`
}

func UserFromLogin(login UserLogin) (User, error) {
	var user User
	if login.UserId == nil || login.Password == nil {
		return user, errors.New("User ID and password have to be provided.")
	}
	user.Id = *login.UserId
	user.SetPassword(*login.Password)
	return user, nil
}

type UserCreate struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
}

func UserFromCreate(create UserCreate) (User, error) {
	var user User
	if create.Name == nil || create.Password == nil {
		return user, errors.New("User name and password have to be provided")
	}
	user.Name = *create.Name
	user.SetPassword(*create.Password)
	return user, nil
}
