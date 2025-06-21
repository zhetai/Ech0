package handler

import "github.com/gin-gonic/gin"

type TodoHandlerInterface interface {
	GetTodoList() gin.HandlerFunc
	AddTodo() gin.HandlerFunc
	UpdateTodo() gin.HandlerFunc
	DeleteTodo() gin.HandlerFunc
}
