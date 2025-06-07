package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/di"
	"github.com/lin-snow/ech0/internal/middleware"
)

type AppRouterGroup struct {
	ResourceGroup     *gin.RouterGroup
	PublicRouterGroup *gin.RouterGroup
	AuthRouterGroup   *gin.RouterGroup
}

func SetupRouter(r *gin.Engine, h *di.Handlers) {
	// Setup Middleware
	setupMiddleware(r)

	// Setup Router Groups
	appRouterGroup := setupRouterGroup(r)

	// Setup Resource Routes
	setupResourceRoutes(appRouterGroup, h)

	// Setup User Routes
	setupUserRoutes(appRouterGroup, h)

	// Setup Echo Routes
	setupEchoRoutes(appRouterGroup, h)

	// Setup Common Routes
	setupCommonRoutes(appRouterGroup, h)

	// Setup Setting Routes
	setupSettingRoutes(appRouterGroup, h)

	// Setup To Do Routes
	setupTodoRoutes(appRouterGroup, h)

	// Setup Connect Routes
	setupConnectRoutes(appRouterGroup, h)

	// Setup No Routes
	setupNoRoutes(r)
}

func setupRouterGroup(r *gin.Engine) *AppRouterGroup {
	resource := r.Group("/")
	public := r.Group("/api")
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuthMiddleware())
	return &AppRouterGroup{
		ResourceGroup:     resource,
		PublicRouterGroup: public,
		AuthRouterGroup:   auth,
	}
}
