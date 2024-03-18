package controllers

import (
	"encoding/json"
	"net/http"
	"todo/db"
	"todo/models"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var userLogin models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user, err := models.UserFromLogin(userLogin)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	token, err := db.LoginUser(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(token)
}
