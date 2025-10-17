package service

import (
	"github.com/gin-gonic/gin"
	model "github.com/lin-snow/ech0/internal/model/metric"
)

type DashboardServiceInterface interface {
	// GetMetrics 获取系统指标
	GetMetrics() (model.Metrics, error)

	// WSSubsribeMetrics 通过 WebSocket 订阅系统指标
	WSSubsribeMetrics(ctx *gin.Context, userId uint) error
}
