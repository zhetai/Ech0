package handler

import "github.com/gin-gonic/gin"

type UserHandlerInterface interface {
	Login() gin.HandlerFunc
	Register() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	UpdateUserAdmin() gin.HandlerFunc
	GetAllUsers() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
	GetUserInfo() gin.HandlerFunc
}
