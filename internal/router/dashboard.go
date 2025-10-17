package router

import "github.com/lin-snow/ech0/internal/di"

func setupDashboardRoutes(appRouterGroup *AppRouterGroup, h *di.Handlers) {
	// Auth
	appRouterGroup.AuthRouterGroup.GET("/dashboard/metrics", h.DashboardHandler.GetMetrics())
	appRouterGroup.AuthRouterGroup.GET("/dashboard/ws/metrics", h.DashboardHandler.WSSubsribeMetrics())
}
