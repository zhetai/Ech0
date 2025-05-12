package routers

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/controllers"
	"github.com/lin-snow/ech0/internal/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 设置 CORS 规则
	r.Use(middleware.Cors())

	// 映射静态文件目录
	r.Use(static.Serve("/", static.LocalFile("./dist", true)))
	r.Static("/api/images", "./data/images")
	r.GET("/rss", controllers.GenerateRSS) // 生成 RSS 订阅链接

	// 公共的路由
	publicRoutes := r.Group("/api")

	publicRoutes.POST("/login", controllers.Login)         // 登录
	publicRoutes.POST("/register", controllers.Register)   // 注册
	publicRoutes.GET("/status", controllers.GetStatus)     // 获取用户信息
	publicRoutes.GET("/settings", controllers.GetSettings) // 获取系统设置
	publicRoutes.GET("/heatmap", controllers.GetHeatMap)   // 获取热力图数据

	// publicRoutes.GET("/messages", controllers.GetMessages) // 获取留言列表

	// 需要鉴权的路由
	authRoutes := r.Group("/api")
	authRoutes.Use(middleware.JWTAuthMiddleware()) // 使用 JWT 鉴权中间件

	// 留言相关路由
	authRoutes.GET("/messages/:id", controllers.GetMessage)          // 获取留言详情
	authRoutes.POST("/messages/page", controllers.GetMessagesByPage) // 分页获取留言列表
	authRoutes.POST("/messages", controllers.PostMessage)            // 发布留言
	authRoutes.DELETE("/messages/:id", controllers.DeleteMessage)    // 删除留言
	authRoutes.POST("/images/upload", controllers.UploadImage)       // 上传图片

	// 用户相关路由
	authRoutes.PUT("/user", controllers.UpdateUser)            // 更新用户信息
	authRoutes.PUT("/user/admin", controllers.UpdateUserAdmin) // 更新用户权限
	authRoutes.GET("/user", controllers.GetUserInfo)           // 获取当前登录的用户信息

	// 设置相关路由
	authRoutes.PUT("/settings", controllers.UpdateSettings) // 更新系统设置

	// 由于Vue3 和SPA模式，所以处理匹配不到的路由(重定向到index.html)
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	return r
}
