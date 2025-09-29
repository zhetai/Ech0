package handler

import "github.com/gin-gonic/gin"

type FediverseHandlerInterface interface {
	// GetActor 获取 Actor 信息
	GetActor(ctx *gin.Context)
}