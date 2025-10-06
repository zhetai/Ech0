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

	// GetS3Settings 获取 S3 存储设置
	GetS3Settings() gin.HandlerFunc

	// UpdateS3Settings 更新 S3 存储设置
	UpdateS3Settings() gin.HandlerFunc

	// GetOAuthSettings 获取 OAuth 设置
	GetOAuthSettings() gin.HandlerFunc

	// UpdateOAuthSettings 更新 OAuth 设置
	UpdateOAuthSettings() gin.HandlerFunc
}
