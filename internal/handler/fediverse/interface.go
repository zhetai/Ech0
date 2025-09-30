package handler

import "github.com/gin-gonic/gin"

type FediverseHandlerInterface interface {
	// Webfinger 处理 Webfinger 请求
	Webfinger(ctx *gin.Context)

	// GetActor 获取 Actor 信息
	GetActor(ctx *gin.Context)

	// PostInbox 处理接收到的 ActivityPub 消息
	PostInbox(ctx *gin.Context)

	// GetOutbox 获取 Outbox 消息
	GetOutbox(ctx *gin.Context)

	// GetFollowers 获取粉丝列表
	GetFollowers(ctx *gin.Context)

	// GetFollowing 获取关注列表
	GetFollowing(ctx *gin.Context)

	// GetObject 获取内容对象
	GetObject(ctx *gin.Context)
}
