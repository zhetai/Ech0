package handler

import "github.com/gin-gonic/gin"

type SettingHandlerInterface interface {
	// GetSettings 获取设置
	GetSettings() gin.HandlerFunc

	// UpdateSettings 更新设置
	UpdateSettings() gin.HandlerFunc
}
