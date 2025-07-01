package handler

import "github.com/gin-gonic/gin"

type ConnectHandlerInterface interface {
	// AddConnect 添加连接
	AddConnect() gin.HandlerFunc

	// DeleteConnect 删除连接
	DeleteConnect() gin.HandlerFunc

	// GetConnectsInfo 获取所有添加的连接的信息
	GetConnectsInfo() gin.HandlerFunc

	// GetConnect 提供当前实例的连接信息
	GetConnect() gin.HandlerFunc

	// GetConnects 获取当前实例添加的所有连接
	GetConnects() gin.HandlerFunc
}
