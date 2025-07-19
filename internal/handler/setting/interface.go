package handler

import "github.com/gin-gonic/gin"

type SettingHandlerInterface interface {
	// GetSettings 获取设置
	GetSettings() gin.HandlerFunc

	// UpdateSettings 更新设置
	UpdateSettings() gin.HandlerFunc

	// GetCommentSettings 获取评论设置
	GetCommentSettings() gin.HandlerFunc

	// UpdateCommentSettings 更新评论设置
	UpdateCommentSettings() gin.HandlerFunc
}
