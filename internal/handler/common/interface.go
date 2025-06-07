package handler

import "github.com/gin-gonic/gin"

type CommonHandlerInterface interface {
	UploadImage(ctx *gin.Context)
	DeleteImage(ctx *gin.Context)
	GetStatus(ctx *gin.Context)
	GetHeatMap(ctx *gin.Context)
	GetRss(ctx *gin.Context)
	UploadAudio(ctx *gin.Context)
	DeleteAudio(ctx *gin.Context)
	GetPlayMusic(ctx *gin.Context)
	PlayMusic(ctx *gin.Context)
}
