package service

import (
	model "github.com/lin-snow/ech0/internal/model/metric"
)

type DashboardServiceInterface interface {
	// GetMetrics 获取系统指标
	GetMetrics() (model.Metrics, error)
}
