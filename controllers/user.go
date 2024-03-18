package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"todo/db"
	"todo/models"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUsers()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userCreate models.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&userCreate); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user, err := models.UserFromCreate(userCreate)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	newUser, err := db.CreateUser(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	token, err := db.LoginUser(newUser)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
