package router

import "github.com/lin-snow/ech0/internal/di"

func setupCommonRoutes(appRouterGroup *AppRouterGroup, h *di.Handlers) {
	// Public
	appRouterGroup.PublicRouterGroup.GET("/status", h.CommonHandler.GetStatus)
	appRouterGroup.PublicRouterGroup.GET("/heatmap", h.CommonHandler.GetHeatMap)
	appRouterGroup.PublicRouterGroup.GET("/getmusic", h.CommonHandler.GetPlayMusic)
	appRouterGroup.PublicRouterGroup.GET("/playmusic", h.CommonHandler.PlayMusic)

	// Auth
	appRouterGroup.AuthRouterGroup.POST("/images/upload", h.CommonHandler.UploadImage)
	appRouterGroup.AuthRouterGroup.POST("/images/delete", h.CommonHandler.DeleteImage)
	appRouterGroup.AuthRouterGroup.POST("/audios/upload", h.CommonHandler.UploadAudio)
	appRouterGroup.AuthRouterGroup.DELETE("/audios/delete", h.CommonHandler.DeleteAudio)
}
