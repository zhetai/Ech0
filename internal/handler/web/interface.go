package handler

import "github.com/gin-gonic/gin"

type WebHandlerInterface interface {
	// Templates 返回前端Web项目编译后的静态文件
	Templates() *gin.HandlerFunc
}
