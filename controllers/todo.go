package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"todo/db"
	"todo/middlewares"
	"todo/models"
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
	userId, err := middlewares.GetUserIdFromAuthenticatedRequest(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todos, err := db.GetTodos(userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	userId, err := middlewares.GetUserIdFromAuthenticatedRequest(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todo.UserId = userId
	todo, err = db.CreateTodo(todo.Title, todo.Text, todo.UserId, todo.IsDone)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userId, err := middlewares.GetUserIdFromAuthenticatedRequest(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todo, err = db.GetTodo(todo.Id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if todo.UserId != userId {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err := db.DeleteTodo(todo.Id); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type TodoUpdateJson struct {
	Id     int               `json:"id"`
	Update models.TodoUpdate `json:"update"`
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userId, err := middlewares.GetUserIdFromAuthenticatedRequest(r)
	if err != nil {
		log.Default().Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var todoUpdateJson TodoUpdateJson
	if err := json.NewDecoder(r.Body).Decode(&todoUpdateJson); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todo, err := db.GetTodo(todoUpdateJson.Id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if todo.UserId != userId {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	todo.Update(todoUpdateJson.Update)
	db.UpdateTodo(todo)
}
