package models

// Todo represents a task or item in a todo list.
// It contains fields such as ID, title, description, user ID, and completion status.
type Todo struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserId int    `json:"userId"`
	IsDone bool   `json:"isDone"`
}

// idea from https://eli.thegreenplace.net/2020/optional-json-fields-in-go/

// TodoUpdate represents optional fields for updating a todo item.
// It allows updating the title, description, and completion status of a todo item.
// Fields with nil values will not be updated.
type TodoUpdate struct {
	Title  *string `json:"title"`
	Text   *string `json:"text"`
	IsDone *bool   `json:"isDone"`
}

// Update updates the fields of a todo item based on the provided TodoUpdate values.
// It updates the title, description, and completion status if non-nil values are provided in TodoUpdate.
// Fields with nil values in TodoUpdate will not be updated.
func (todo *Todo) Update(valuesFrom TodoUpdate) {
	if valuesFrom.Title != nil {
		todo.Title = *valuesFrom.Title
	}
	if valuesFrom.Text != nil {
		todo.Text = *valuesFrom.Text
	}
	if valuesFrom.IsDone != nil {
		todo.IsDone = *valuesFrom.IsDone
	}
}
