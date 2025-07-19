package router

import "github.com/lin-snow/ech0/internal/di"

// setupSettingRoutes 设置设置路由
func setupSettingRoutes(appRouterGroup *AppRouterGroup, h *di.Handlers) {
	// Public
	appRouterGroup.PublicRouterGroup.GET("/settings", h.SettingHandler.GetSettings())
	appRouterGroup.PublicRouterGroup.GET("/comment/settings", h.SettingHandler.GetCommentSettings())

	// Auth
	appRouterGroup.AuthRouterGroup.PUT("/settings", h.SettingHandler.UpdateSettings())
	appRouterGroup.AuthRouterGroup.PUT("/comment/settings", h.SettingHandler.UpdateCommentSettings())
}
