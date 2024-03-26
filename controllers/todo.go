package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"todo/db"
	"todo/logger"
	"todo/middlewares"
	"todo/models"
)

func getTodoFromPathId(r *http.Request, w http.ResponseWriter) (models.Todo, error) {
	var todo models.Todo
	todoId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "No todo id was given in the request path", http.StatusBadRequest)
		return todo, errors.New("No proper id was given for a todo in the request path")
	}
	todo, err = db.GetTodo(todoId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return todo, err
	}
	return todo, nil
}

// GetTodo returns a [models.Todo] with an ID specified as a request path value. It also requires that the user is
// authorized by [middlewares.AuthenticateUser] as it expects the request's context to have a user object.
// If the query for the [models.Todo] succeeds it is returned as JSON in the response body.
func GetTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	todo, err := getTodoFromPathId(r, w)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if todo.UserId != user.Id {
		//check if user is todoUser or has share on todo, then return
		userShares, err := db.GetTodoShares(todo.Id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !slices.Contains(userShares, user.Id) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	}
	json.NewEncoder(w).Encode(todo)
}

// GetTodos returns a list of [models.Todo] with an ID specified as a request path value. It also requires that the user is
// authorized by [middlewares.AuthenticateUser] as it expects the request's context to have a user object.
// If the query for the list succeeds it is returned as JSON in the response body.
func GetTodos(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
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

// CreateTodo creates a [models.Todo] based on the corresponding fields in the request body. It expects the request
// to be authorized by [middlewares.AuthenticateUser] as it expects the request's context to have a user object.
// If the the object is created and persisted successfully it is returned as JSON in the response body.
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

// DeleteTodo deletes a [models.Todo] based on the request's path value. It expects the request
// to be authorized by [middlewares.AuthenticateUser] as it expects the request's context to have a user object.
// If the the object is deleted successfully it is returned as JSON in the response body.
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	todo, err := getTodoFromPathId(r, w)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if todo.UserId != user.Id {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	deletedTodo, err := db.DeleteTodo(todo.Id)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(deletedTodo)
}

// UpdateTodo updates a [models.Todo] based on the corresponding fields in the request body. The [models.Todo] to update is
// specified as a path value. The request body with the fields to update must not provide all fields of a [models.Todo]
// but only the fields that should be updated. It expects the request to be authorized by [middlewares.AuthenticateUser]
// as it expects the request's context to have a user object. If the the object is updated successfully the updated
// object is returned as JSON in the response body.
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	todo, err := getTodoFromPathId(r, w)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if todo.UserId != user.Id {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	var todoUpdate models.TodoUpdate
	if err := json.NewDecoder(r.Body).Decode(&todoUpdate); err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todo.Update(todoUpdate)
	if err := db.UpdateTodo(todo); err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

type TodoShare struct {
	Id     int  `json:"id"`
	TodoId *int `json:"todoId"`
	UserId *int `json:"userId"`
}

// ShareTodo shares a [model.Todo] with another [models.User] that is not the creator of the [models.Todo].
func ShareTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	todo, err := getTodoFromPathId(r, w)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	var todoShare TodoShare
	if err := json.NewDecoder(r.Body).Decode(&todoShare); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todoShare.TodoId = &todo.Id
	if todoShare.TodoId == nil || todoShare.UserId == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if user.Id != todo.UserId {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if todo.UserId == *todoShare.UserId {
		http.Error(w, "You are trying to share your own Todo with yourself.", http.StatusBadRequest)
		return
	}
	shareId, err := db.CreateTodoShare(*todoShare.TodoId, *todoShare.UserId)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	todoShare.Id = shareId
	if err := json.NewEncoder(w).Encode(todoShare); err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// UnshareTodo removes a [models.User] from the users a [model.Todo] is shared with.
func UnshareTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.ContextUserKey).(models.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	todo, err := getTodoFromPathId(r, w)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	var todoShare TodoShare
	if err := json.NewDecoder(r.Body).Decode(&todoShare); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	todoShare.TodoId = &todo.Id
	if todoShare.TodoId == nil || todoShare.UserId == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if user.Id != todo.UserId {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	shareId, err := db.DeleteTodoShare(*todoShare.TodoId, *todoShare.UserId)
	todoShare.Id = shareId
	if err := json.NewEncoder(w).Encode(todoShare); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
