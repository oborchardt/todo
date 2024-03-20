package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"todo/db"
)

const ContextUserKey = "user"

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

func AuthenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := preprocessToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := db.AuthenticateUser(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
