package service

import (
	model "github.com/lin-snow/ech0/internal/model/metric"
	"github.com/lin-snow/ech0/internal/monitor"
)

type DashboardService struct {
	monitor *monitor.Monitor
}

func NewDashboardService(monitor *monitor.Monitor) DashboardServiceInterface {
	return &DashboardService{
		monitor: monitor,
	}
}

func (s *DashboardService) GetMetrics() (model.Metrics, error) {
	return s.monitor.GetMetrics(), nil
}
