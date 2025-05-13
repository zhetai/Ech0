package services

import (
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
)

func GetTodos(userID uint) ([]models.Todo, error) {
	todos, err := repository.GetTodosByUserID(userID)
	if err != nil {
		return todos, err
	}

	// 除去已完成的 Todo
	for i := len(todos) - 1; i >= 0; i-- {
		if todos[i].Status == models.NotDone {
			todos = append(todos[:i], todos[i+1:]...)
		}
	}

	return todos, nil
}

func AddTodo(todo *models.Todo) error {
	todo.Status = models.NotDone

	if err := repository.CreateTodo(todo); err != nil {
		return err
	}
	return nil
}

func UpdateTodo(todo *models.Todo) error {
	// 设置Todo的状态
	if todo.Status == models.NotDone {
		todo.Status = models.Done
	} else {
		todo.Status = models.NotDone
	}

	if err := repository.UpdateTodo(todo); err != nil {
		return err
	}

	return nil
}

func DeleteTodo(todo *models.Todo) error {
	if err := repository.DeleteTodo(todo.ID); err != nil {
		return err
	}

	return nil
}
