package handler

import "github.com/gin-gonic/gin"

type CommonHandlerInterface interface {
	UploadImage() gin.HandlerFunc
	DeleteImage() gin.HandlerFunc
	GetStatus() gin.HandlerFunc
	GetHeatMap() gin.HandlerFunc
	UploadAudio() gin.HandlerFunc
	DeleteAudio() gin.HandlerFunc
	GetPlayMusic() gin.HandlerFunc
	GetRss(ctx *gin.Context)
	PlayMusic(ctx *gin.Context)
}
