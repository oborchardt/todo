package db

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/mattn/go-sqlite3"
	"time"
	"todo/logger"
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
	// while loop if random token is not unique
	created := false
	var token string
	for !created {
		var err error
		token, err = createUserToken(tokenLength)
		if err != nil {
			return "", err
		}
		if err = getDb().QueryRow(stmt, token, expiration, userId).Scan(&token); err != nil {
			var sqliteErr sqlite3.Error
			if errors.As(err, &sqliteErr) {
				if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					logger.Warning(sqliteErr.Error())
					continue
				}
			}
			return "", err
		}
		created = true
	}

	return token, nil
}

func LoginUser(user models.User) (string, error) {
	// check user and password
	var password string
	stmt := `SELECT id, password FROM users WHERE name = ?`
	if err := getDb().QueryRow(stmt, user.Name).Scan(&user.Id, &password); err != nil {
		return "", err
	}
	if err := user.CheckPassword([]byte(password)); err != nil {
		return "", err
	}
	// create a token for the user
	return updateToken(user.Id)
}

func AuthenticateUser(token string) (models.User, error) {
	var user models.User
	expiration := time.Now()
	stmt := `SELECT id, name FROM users WHERE token = ? AND expiration > ?`
	err := getDb().QueryRow(stmt, token, expiration).Scan(&user.Id, &user.Name)
	return user, err
}
