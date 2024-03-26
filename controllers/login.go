package controllers

import (
	"encoding/json"
	"net/http"
	"todo/db"
	"todo/logger"
	"todo/models"
)

type loginResponse struct {
	Token string `json:"token"`
}

// LoginUser uses a username and a password provided in the request body to authenticate a user.
// It returns a JSON object that contains a token in the response body.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var userLogin models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user, err := models.UserFromLogin(userLogin)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	token, err := db.LoginUser(user)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}
