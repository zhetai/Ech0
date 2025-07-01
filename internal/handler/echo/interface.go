package handler

import "github.com/gin-gonic/gin"

type EchoHandlerInterface interface {
	// PostEcho 发布新的 Echo
	PostEcho() gin.HandlerFunc

	// GetEchosByPage 获取 Echo 列表，支持分页
	GetEchosByPage() gin.HandlerFunc

	// DeleteEcho 删除 Echo
	DeleteEcho() gin.HandlerFunc

	// GetTodayEchos 获取今天的 Echo 列表
	GetTodayEchos() gin.HandlerFunc

	// UpdateEcho 更新 Echo
	UpdateEcho() gin.HandlerFunc
}
