package router

import "github.com/lin-snow/ech0/internal/di"

func setupDashboardRoutes(appRouterGroup *AppRouterGroup, h *di.Handlers) {
	// Auth
	appRouterGroup.AuthRouterGroup.GET("/dashboard/metrics", h.DashboardHandler.GetMetrics())
	appRouterGroup.WSRouterGroup.GET("/dashboard/metrics", h.DashboardHandler.WSSubsribeMetrics())
}
