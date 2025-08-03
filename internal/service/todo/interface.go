package service

import (
	model "github.com/lin-snow/ech0/internal/model/todo"
)

type TodoServiceInterface interface {
	// GetTodoList 获取当前用户的 To do 列表
	GetTodoList(userid uint) ([]model.Todo, error)

	// AddTodo 创建新的 To do
	AddTodo(userid uint, todo *model.Todo) error

	// UpdateTodo 更新指定ID的 To do
	UpdateTodo(userid uint, id int64) error

	// DeleteTodo 删除指定ID的 To do
	DeleteTodo(userid uint, id int64) error
}
