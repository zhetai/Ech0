package handler

import "github.com/gin-gonic/gin"

type ConnectHandlerInterface interface {
	AddConnect() gin.HandlerFunc
	DeleteConnect() gin.HandlerFunc
	GetConnectsInfo() gin.HandlerFunc
	GetConnect() gin.HandlerFunc
	GetConnects() gin.HandlerFunc
}
