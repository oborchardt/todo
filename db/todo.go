package db

import (
	"database/sql"
	"errors"
	"todo/models"
)

func CreateTodo(title string, text string, userId int, isDone bool) (models.Todo, error) {
	stmt := `INSERT INTO todos(title, text, user_id, is_done) VALUES (?, ?, ?, ?) RETURNING id`
	id := 0
	err := getDb().QueryRow(stmt, title, text, userId, false).Scan(&id)
	var todo models.Todo
	if err != nil {
		return todo, err
	}
	todo = models.Todo{
		Id:     id,
		Title:  title,
		Text:   text,
		UserId: userId,
		IsDone: isDone,
	}
	return todo, nil
}

func GetTodos(userId int, includeShared bool) ([]models.Todo, error) {
	var todos []models.Todo
	var stmt string
	if includeShared {
		stmt = `SELECT t.id, t.title, t.text, t.is_done, t.user_id FROM todos AS t LEFT OUTER JOIN users_todos AS u ON (t.id = u.todo_id) WHERE t.user_id = $1 OR u.user_id = $1`
	} else {
		stmt = `SELECT id, title, text, is_done, user_id FROM todos WHERE user_id = ?`
	}
	rows, err := getDb().Query(stmt, userId)
	if err != nil {
		return todos, err
	}
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Text, &todo.IsDone, &todo.UserId)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func GetTodo(todoId int) (models.Todo, error) {
	var todo models.Todo
	stmt := `SELECT id, title, text, is_done, user_id FROM todos WHERE id = ?`
	err := getDb().QueryRow(stmt, todoId).Scan(&todo.Id, &todo.Title, &todo.Text, &todo.IsDone, &todo.UserId)
	return todo, err
}

func UpdateTodo(todo models.Todo) error {
	stmt := `UPDATE todos SET title = ?, text = ?, is_done = ?, user_id = ? WHERE id = ?`
	err := getDb().QueryRow(stmt, todo.Title, todo.Text, todo.IsDone, todo.UserId, todo.Id).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func DeleteTodo(todoId int) error {
	stmt := `DELETE FROM todos WHERE id = ?;`
	err := getDb().QueryRow(stmt, todoId).Scan()
	// query returns no rows but that is intended
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func CreateTodoShare(todoId int, userId int) (int, error) {
	stmt := `INSERT INTO users_todos(todo_id, user_id) VALUES (?, ?) RETURNING id`
	var shareId int
	err := getDb().QueryRow(stmt, todoId, userId).Scan(&shareId)
	return shareId, err
}

func DeleteTodoShare(todoId int, userId int) (int, error) {
	stmt := `DELETE FROM users_todos WHERE todo_id = ? AND user_id = ? RETURNING id`
	var shareId int
	err := getDb().QueryRow(stmt, todoId, userId).Scan(&shareId)
	return shareId, err
}
