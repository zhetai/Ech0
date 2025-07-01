package handler

import "github.com/gin-gonic/gin"

type TodoHandlerInterface interface {
	// GetTodoList 获取待办事项列表
	GetTodoList() gin.HandlerFunc

	// AddTodo 添加新的待办事项
	AddTodo() gin.HandlerFunc

	// UpdateTodo 更新待办事项
	UpdateTodo() gin.HandlerFunc

	// DeleteTodo 删除待办事项
	DeleteTodo() gin.HandlerFunc
}
