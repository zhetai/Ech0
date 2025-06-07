package router

import (
	"github.com/lin-snow/ech0/internal/di"
)

func setupResourceRoutes(appRouterGroup *AppRouterGroup, h *di.Handlers) {
	// TODO: 不知道为什么放在这里会失效，无法使用
	//appRouterGroup.ResourceGroup.Use(static.Serve("/", static.LocalFile("./template", true)))
	//appRouterGroup.ResourceGroup.Static("/app.webmanifest", "./template/app.webmanifest")
	//appRouterGroup.ResourceGroup.Static("api/images", "./data/images")

	appRouterGroup.ResourceGroup.GET("rss", h.CommonHandler.GetRss)
}
