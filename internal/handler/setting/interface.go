package handler

import "github.com/gin-gonic/gin"

type SettingHandlerInterface interface {
	GetSettings() gin.HandlerFunc
	UpdateSettings() gin.HandlerFunc
}
