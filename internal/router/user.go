package router

import "github.com/lin-snow/ech0/internal/di"

// setupUserRoutes 设置用户路由
func setupUserRoutes(appRouterGroup *AppRouterGroup, h *di.Handlers) {
	// OAuth2
	appRouterGroup.ResourceGroup.GET("/oauth/github/login", h.UserHandler.GitHubLogin())
	appRouterGroup.ResourceGroup.GET("/oauth/github/callback", h.UserHandler.GitHubCallback())

	// Public
	appRouterGroup.PublicRouterGroup.POST("/login", h.UserHandler.Login())
	appRouterGroup.PublicRouterGroup.POST("/register", h.UserHandler.Register())
	appRouterGroup.PublicRouterGroup.GET("/allusers", h.UserHandler.GetAllUsers())

	// Auth
	appRouterGroup.AuthRouterGroup.GET("/user", h.UserHandler.GetUserInfo())
	appRouterGroup.AuthRouterGroup.PUT("/user", h.UserHandler.UpdateUser())
	appRouterGroup.AuthRouterGroup.DELETE("/user/:id", h.UserHandler.DeleteUser())
	appRouterGroup.AuthRouterGroup.PUT("/user/admin/:id", h.UserHandler.UpdateUserAdmin())
}
