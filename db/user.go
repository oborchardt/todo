package db

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"time"
	"todo/logger"
	"todo/models"
)

// CreateUser inserts a [models.User] into the database. On success the user is returned.
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

// GetUsers retrieves a list of [models.User] from the database.
// It executes a SQL query to fetch the user IDs and names from the 'users' table.
// The function returns a slice of models.User containing the retrieved users and an error.
// If the query execution encounters an error, it returns an empty slice of users and the error.
// If the query executes successfully, it iterates over the result set, populates the users,
// and appends them to the slice. Finally, it returns the populated slice of users and a nil error.
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
					newErr := fmt.Errorf("updateToken: Token generation collision: %w", err)
					logger.Warning(newErr.Error())
					continue
				}
			}
			return "", err
		}
		created = true
	}

	return token, nil
}

// LoginUser verifies the credentials of a user attempting to log in and generates an authentication token upon successful validation.
// It takes a [models.User] object representing the user attempting to log in.
// The function executes a SQL query to retrieve the user's ID and hashed password from the 'users' table based on the provided username.
// If the user is found and the password matches, it generates an authentication token for the user and returns it along with a nil error.
// If the provided username is not found in the database or the password doesn't match, it returns an empty string and an error.
// If any database operation fails, it returns an error wrapping the original error encountered during the database interaction.
func LoginUser(user models.User) (string, error) {
	// check user and password
	var userId int
	var password string
	stmt := `SELECT id, password FROM users WHERE name = ?`
	if err := getDb().QueryRow(stmt, user.Name).Scan(&userId, &password); err != nil {
		return "", fmt.Errorf("LoginUser: failed to find user %s: %w", user.Name, err)
	}
	if err := user.CheckPassword([]byte(password)); err != nil {
		return "", err
	}
	// create a token for the user
	return updateToken(userId)
}

// AuthenticateUser verifies the validity of an authentication token and retrieves the corresponding user.
// It takes a token string as input representing the authentication token to be verified.
// The function executes a SQL query to retrieve the user's ID and name from the 'users' table based on the provided token.
// It also ensures that the token has not expired by checking the expiration time against the current time.
// If the token is valid and corresponds to an existing user whose token has not expired,
// it returns the user object along with a nil error.
// If the token is invalid, expired, or if any database operation fails, it returns an empty user object and an error.
func AuthenticateUser(token string) (models.User, error) {
	var user models.User
	expiration := time.Now()
	stmt := `SELECT id, name FROM users WHERE token = ? AND expiration > ?`
	err := getDb().QueryRow(stmt, token, expiration).Scan(&user.Id, &user.Name)
	return user, err
}
