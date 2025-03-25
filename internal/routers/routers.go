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
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	r.Static("/api/images", "./data/images")

	r.GET("/rss", controllers.GenerateRSS) // 生成 RSS 订阅链接

	// 公共的路由
	publicRoutes := r.Group("/api")

	publicRoutes.POST("/login", controllers.Login)       // 登录
	publicRoutes.POST("/register", controllers.Register) // 注册
	publicRoutes.GET("/status", controllers.GetStatus)   // 获取用户信息

	// publicRoutes.GET("/messages", controllers.GetMessages) // 获取留言列表

	// 需要鉴权的路由
	authRoutes := r.Group("/api")
	authRoutes.Use(middleware.JWTAuthMiddleware()) // 使用 JWT 鉴权中间件

	authRoutes.GET("/messages/:id", controllers.GetMessage)          // 获取留言详情
	authRoutes.POST("/messages/page", controllers.GetMessagesByPage) // 分页获取留言列表

	authRoutes.POST("/messages", controllers.PostMessage)         // 发布留言
	authRoutes.DELETE("/messages/:id", controllers.DeleteMessage) // 删除留言

	// 添加图片上传路由
	authRoutes.POST("/images/upload", controllers.UploadImage) // 上传图片

	// 更新用户信息
	authRoutes.PUT("/user/change_password", controllers.ChangePassword) // 修改密码
	authRoutes.PUT("/user/update", controllers.UpdateUser)              // 更新用户信息
	authRoutes.PUT("/user/admin", controllers.UpdateUserAdmin)          // 更新用户权限

	// 获取当前登录的用户信息
	authRoutes.GET("/user", controllers.GetUserInfo) // 获取当前登录的用户信息

	// 更新系统设置
	// authRoutes.PUT("/setting/update", controllers.UpdateSetting) // 更新系统设置

	return r
}
