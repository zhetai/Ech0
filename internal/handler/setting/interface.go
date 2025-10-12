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

	// GetOAuth2Settings 获取 OAuth2 设置
	GetOAuth2Settings() gin.HandlerFunc

	// UpdateOAuth2Settings 更新 OAuth2 设置
	UpdateOAuth2Settings() gin.HandlerFunc

	// GetOAuth2Status 获取 OAuth2 状态
	GetOAuth2Status() gin.HandlerFunc

	// GetWebhook 获取所有 Webhook
	GetWebhook() gin.HandlerFunc

	// DeleteWebhook 删除 Webhook
	DeleteWebhook() gin.HandlerFunc

	// UpdateWebhook 更新 Webhook
	UpdateWebhook() gin.HandlerFunc

	// CreateWebhook 创建 Webhook
	CreateWebhook() gin.HandlerFunc
}
