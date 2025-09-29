package handler

import "github.com/gin-gonic/gin"

type FediverseHandlerInterface interface {
	// GetActor 获取 Actor 信息
	GetActor(ctx *gin.Context)

	// PostInbox 处理接收到的 ActivityPub 消息
	PostInbox(ctx *gin.Context)
}