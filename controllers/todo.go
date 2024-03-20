package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo/db"
	"todo/middlewares"
	"todo/models"
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	fmt.Println(user)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	//check if shared flag is set
	includeShared := r.URL.Query().Has("shared")
	todos, err := db.GetTodos(user.Id, includeShared)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todo.UserId = user.Id
	todo, err := db.CreateTodo(todo.Title, todo.Text, todo.UserId, todo.IsDone)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todo, err := db.GetTodo(todo.Id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if todo.UserId != user.Id {
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
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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
	if todo.UserId != user.Id {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	todo.Update(todoUpdateJson.Update)
	db.UpdateTodo(todo)
}

type TodoShare struct {
	Id     int  `json:"id"`
	TodoId *int `json:"todoId"`
	UserId *int `json:"userId"`
}

// CanEditShare(user_id, todo_id / todo_object) funktion o.ä um Code-Duplizierung zu vermeiden!
func ShareTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var todoShare TodoShare
	if err := json.NewDecoder(r.Body).Decode(&todoShare); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if todoShare.TodoId == nil || todoShare.UserId == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todo, err := db.GetTodo(*todoShare.TodoId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if user.Id != todo.UserId {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	shareId, err := db.CreateTodoShare(*todoShare.TodoId, *todoShare.UserId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	todoShare.Id = shareId
	if err := json.NewEncoder(w).Encode(todoShare); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func UnshareTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var todoShare TodoShare
	if err := json.NewDecoder(r.Body).Decode(&todoShare); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if todoShare.TodoId == nil || todoShare.UserId == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todo, err := db.GetTodo(*todoShare.TodoId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if user.Id != todo.UserId {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	shareId, err := db.DeleteTodoShare(*todoShare.TodoId, *todoShare.UserId)
	todoShare.Id = shareId
	if err := json.NewEncoder(w).Encode(todoShare); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
