package routers

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/controllers"
	"github.com/lin-snow/ech0/internal/middleware"
)

func SetupRouter(r *gin.Engine) {
	// 设置 CORS 规则
	r.Use(middleware.Cors())

	// 映射静态文件目录
	r.Use(static.Serve("/", static.LocalFile("./template", true)))
	r.Static("/api/images", "./data/images")

	// 处理前端路由
	r.GET("/rss", controllers.GenerateRSS) // 生成 RSS 订阅链接

	// 设置公共路由
	setupPublicRoutes(r)

	// 设置需要鉴权的路由
	setupAuthRoutes(r)

	// 由于Vue3 和SPA模式，所以处理匹配不到的路由(重定向到index.html)
	r.NoRoute(func(c *gin.Context) {
		c.File("./template/index.html")
	})
}
