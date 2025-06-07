package router

import (
	"github.com/gin-contrib/static"
	"github.com/lin-snow/ech0/internal/di"
)

func setupResourceRoutes(appRouterGroup *AppRouterGroup, h *di.Handlers) {
	appRouterGroup.ResourceGroup.Use(static.Serve("/", static.LocalFile("./template", true)))
	appRouterGroup.ResourceGroup.Static("api/images", "./data/images")
	appRouterGroup.ResourceGroup.GET("rss", h.CommonHandler.GetRss)
}
