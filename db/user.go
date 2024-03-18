package db

import (
	"crypto/rand"
	"encoding/hex"
	"time"
	"todo/models"
)

func CreateUser(user models.User) (models.User, error) {
	stmt := `INSERT INTO users (name, password) VALUES (?, ?) RETURNING id;`
	id := 0
	pw, err := user.GetPasswordHash()
	if err != nil {
		return user, err
	}
	err = getDb().QueryRow(stmt, user.Name, pw).Scan(&id)
	if err != nil {
		return user, err
	}
	user.Id = id
	return user, nil
}

func GetUsers() ([]models.User, error) {
	var users []models.User
	stmt := `SELECT id, name FROM users`
	rows, err := getDb().Query(stmt)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func createUserToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func updateToken(userId int) (string, error) {
	tokenLength := 32
	expiration := time.Now().Add(time.Minute * time.Duration(5))
	stmt := `UPDATE users SET token = ?, expiration = ? WHERE id = ? RETURNING token`
	// add while loop if random token is not unique
	token, err := createUserToken(tokenLength)
	if err != nil {
		return "", err
	}
	if err = getDb().QueryRow(stmt, token, expiration, userId).Scan(&token); err != nil {
		return "", err
	}
	return token, nil
}

func LoginUser(user models.User) (string, error) {
	// check user and password
	var password string
	stmt := `SELECT password FROM users WHERE id = ?`
	if err := getDb().QueryRow(stmt, user.Id).Scan(&password); err != nil {
		return "", err
	}
	if err := user.CheckPassword([]byte(password)); err != nil {
		return "", err
	}
	// create a token for the user
	return updateToken(user.Id)
}

func AuthenticateUser(token string) (int, error) {
	var userId int
	expiration := time.Now()
	stmt := `SELECT id FROM users WHERE token = ? AND expiration > ?`
	err := getDb().QueryRow(stmt, token, expiration).Scan(&userId)
	return userId, err
}
