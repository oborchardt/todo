package middlewares

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"todo/db"
)

func preprocessToken(token string) (string, error) {
	split := strings.Split(token, " ")
	if len(split) != 2 {
		return "", errors.New("malformatted token")
	}
	if split[0] != "Bearer" {
		return "", errors.New("wrong authorization method")
	}
	return split[1], nil
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := preprocessToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		userId, err := db.AuthenticateUser(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		r.Header.Set("UserId", strconv.Itoa(userId))
		next.ServeHTTP(w, r)
	})
}

func GetUserIdFromAuthenticatedRequest(r *http.Request) (int, error) {
	return strconv.Atoi(r.Header.Get("UserId"))
}
