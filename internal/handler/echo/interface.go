package handler

import "github.com/gin-gonic/gin"

type EchoHandlerInterface interface {
	PostEcho() gin.HandlerFunc
	GetEchosByPage() gin.HandlerFunc
	DeleteEcho() gin.HandlerFunc
	GetTodayEchos() gin.HandlerFunc
}
