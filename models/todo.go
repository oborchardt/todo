package models

type Todo struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserId int    `json:"userId"`
	IsDone bool   `json:"isDone"`
}

// idea from https://eli.thegreenplace.net/2020/optional-json-fields-in-go/
type TodoUpdate struct {
	Title  *string `json:"title"`
	Text   *string `json:"text"`
	IsDone *bool   `json:"isDone"`
}

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
