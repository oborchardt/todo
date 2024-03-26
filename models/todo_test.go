package models

import "testing"

func TestTodo_Update(t *testing.T) {
	todo := Todo{
		Id:     1,
		Title:  "Test 1",
		Text:   "test 1",
		UserId: 3,
		IsDone: false,
	}
	updatedTodo := Todo{
		Id:     1,
		Title:  "TEST 2",
		Text:   "test 2",
		UserId: 3,
		IsDone: true,
	}
	todoUpdate := TodoUpdate{
		Title:  &updatedTodo.Title,
		Text:   &updatedTodo.Text,
		IsDone: &updatedTodo.IsDone,
	}
	todo.Update(todoUpdate)
	if todo != updatedTodo {
		t.Fatalf("Todo update went wrong.")
	}
	todo.Update(TodoUpdate{
		Title:  nil,
		Text:   nil,
		IsDone: nil,
	})
	if todo != updatedTodo {
		t.Fatalf("Fields were changed but they shouldn't have changed.")
	}
}
