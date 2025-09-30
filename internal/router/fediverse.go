package router

import "github.com/lin-snow/ech0/internal/di"

func setupFediverseRoutes(appRouterGroup *AppRouterGroup, h *di.Handlers) {
	// // ActivityPub discovery
	appRouterGroup.PublicRouterGroup.GET("/.well-known/webfinger", h.FediverseHandler.Webfinger)

	// Actor (用户资料)
	appRouterGroup.ResourceGroup.GET("/users/:username", h.FediverseHandler.GetActor)

	// Inbox (接收 ActivityPub 消息)
	appRouterGroup.PublicRouterGroup.POST("/users/:username/inbox", h.FediverseHandler.PostInbox)

	// Outbox (发布消息)
	appRouterGroup.PublicRouterGroup.GET("/users/:username/outbox", h.FediverseHandler.GetOutbox)

	// Followers list
	appRouterGroup.PublicRouterGroup.GET("/users/:username/followers", h.FediverseHandler.GetFollowers)

	// Following list
	appRouterGroup.PublicRouterGroup.GET("/users/:username/following", h.FediverseHandler.GetFollowing)

	// Objects (内容对象访问)
	appRouterGroup.PublicRouterGroup.GET("/objects/:id", h.FediverseHandler.GetObject)
}
