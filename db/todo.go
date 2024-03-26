package db

import (
	"database/sql"
	"errors"
	"todo/models"
)

// CreateTodo creates a [models.Todo] from the parameters and inserts it into the database. If the insert was successful
// the object is returned, otherwise an error is returned
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

// GetTodos fetches the [models.Todo]s created by a [models.User] with id userId. If includeShared is set to true, all
// the todos shared with that user are also included
func GetTodos(userId int, includeShared bool) ([]models.Todo, error) {
	var todos []models.Todo
	var stmt string
	if includeShared {
		stmt = `SELECT DISTINCT(t.id), t.title, t.text, t.is_done, t.user_id FROM todos AS t LEFT JOIN users_todos AS u ON (t.id = u.todo_id) WHERE t.user_id = $1 OR u.user_id = $1`
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

// GetTodo fetches the [models.Todo] with the [models.Todo.Id] equal to todoId from the database.
func GetTodo(todoId int) (models.Todo, error) {
	var todo models.Todo
	stmt := `SELECT id, title, text, is_done, user_id FROM todos WHERE id = ?`
	err := getDb().QueryRow(stmt, todoId).Scan(&todo.Id, &todo.Title, &todo.Text, &todo.IsDone, &todo.UserId)
	return todo, err
}

// UpdateTodo sets all values of the row in the database equal to the editable fields of todo. The id for the row is
// taken from the [models.Todo.Id] of todo.
func UpdateTodo(todo models.Todo) error {
	stmt := `UPDATE todos SET title = ?, text = ?, is_done = ?, user_id = ? WHERE id = ?`
	err := getDb().QueryRow(stmt, todo.Title, todo.Text, todo.IsDone, todo.UserId, todo.Id).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

// DeleteTodo deletes a [models.Todo] from the database where its [models.Todo.Id] is equal to todoId
func DeleteTodo(todoId int) (models.Todo, error) {
	var todo models.Todo
	stmt := `DELETE FROM todos WHERE id = ? RETURNING id, title, text, is_done, user_id;`
	err := getDb().QueryRow(stmt, todoId).Scan(&todo.Id, &todo.Title, &todo.Text, &todo.IsDone, &todo.UserId)
	return todo, err
}

// GetTodoShares fetches all the user ids the [models.Todo] is shared with
func GetTodoShares(todoId int) ([]int, error) {
	var shares []int
	stmt := `SELECT user_id FROM users_todos WHERE todo_id = ?`
	rows, err := getDb().Query(stmt, todoId)
	if err != nil {
		return shares, err
	}
	for rows.Next() {
		var share int
		err := rows.Scan(&share)
		if err != nil {
			return shares, err
		}
		shares = append(shares, share)
	}
	return shares, nil
}

// CreateTodoShare inserts a share of a [models.Todo] with user userId in the database
func CreateTodoShare(todoId int, userId int) (int, error) {
	stmt := `INSERT INTO users_todos(todo_id, user_id) VALUES (?, ?) RETURNING id`
	var shareId int
	err := getDb().QueryRow(stmt, todoId, userId).Scan(&shareId)
	return shareId, err
}

// CreateTodoShare deletes a share of a [models.Todo] with user userId in the database
func DeleteTodoShare(todoId int, userId int) (int, error) {
	stmt := `DELETE FROM users_todos WHERE todo_id = ? AND user_id = ? RETURNING id`
	var shareId int
	err := getDb().QueryRow(stmt, todoId, userId).Scan(&shareId)
	return shareId, err
}
