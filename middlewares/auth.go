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

// AuthenticateUser is a middleware that authenticates incoming requests by verifying the provided authorization token.
// It takes an [http.HandlerFunc] 'next' as input, representing the next HTTP handler function to be executed in the chain.
// The middleware intercepts the incoming request, extracts the authorization token from the request header,
// and preprocesses it to remove any prefix or formatting.
// It then attempts to authenticate the user based on the extracted token by calling the [db.AuthenticateUser] function.
// If the token is valid and corresponds to an existing user, it adds the authenticated user to the request context
// using the ContextUserKey key.
// If the token is invalid or authentication fails, it responds with an HTTP status code 401 (Unauthorized).
// If the authorization token is missing or malformed, it responds with an HTTP status code 400 (Bad Request).
// After authentication, the middleware passes the request to the next handler in the chain.
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
