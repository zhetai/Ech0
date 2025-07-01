package router

import (
	"github.com/gin-gonic/gin"
)

// setupNoRoutes 设置无路由处理
func setupNoRoutes(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		c.File("./template/index.html")
	})
}
