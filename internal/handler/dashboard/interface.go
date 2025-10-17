package handler

import "github.com/gin-gonic/gin"

type DashboardHandlerInterface interface {
	// GetMetrics 获取系统指标
	GetMetrics() gin.HandlerFunc

	// WSSubsribeMetrics 通过 WebSocket 订阅系统指标
	WSSubsribeMetrics() gin.HandlerFunc
}
