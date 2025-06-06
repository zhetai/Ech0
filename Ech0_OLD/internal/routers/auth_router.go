package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/controllers"
	"github.com/lin-snow/ech0/internal/middleware"
)

func setupAuthRoutes(r *gin.Engine) {
	// 需要鉴权的路由
	authRoutes := r.Group("/api")
	authRoutes.Use(middleware.JWTAuthMiddleware()) // 使用 JWT 鉴权中间件

	setupMessageRoutes(authRoutes)  // 留言相关路由
	setupTodoRoutes(authRoutes)     // Todo 相关路由
	setupMusicRoutes(authRoutes)    // 音乐相关路由
	setupUserRoutes(authRoutes)     // 用户相关路由
	setupSettingsRoutes(authRoutes) // 设置相关路由
	setupConnectRoutes(authRoutes)  // Connect相关路由
}

// 留言相关路由
func setupMessageRoutes(rg *gin.RouterGroup) {
	rg.GET("/messages/:id", controllers.GetMessage)          // 获取留言详情
	rg.POST("/messages/page", controllers.GetMessagesByPage) // 分页获取留言列表
	rg.POST("/messages", controllers.PostMessage)            // 发布留言
	rg.DELETE("/messages/:id", controllers.DeleteMessage)    // 删除留言
	rg.POST("/images/upload", controllers.UploadImage)       // 上传图片
	rg.DELETE("/images/delete", controllers.DeleteImage)     // 删除图片
}

// Todo 相关路由
func setupTodoRoutes(rg *gin.RouterGroup) {
	rg.GET("/todo", controllers.GetTodos)          // 获取 Todo 列表
	rg.POST("/todo", controllers.PostTodo)         // 发布 Todo
	rg.PUT("/todo/:id", controllers.UpdateTodo)    // 更新 Todo
	rg.DELETE("/todo/:id", controllers.DeleteTodo) // 删除 Todo
}

// 音乐相关路由
func setupMusicRoutes(rg *gin.RouterGroup) {
	rg.POST("/audios/upload", controllers.UploadAudio)   // 上传音频
	rg.DELETE("/audios/delete", controllers.DeleteAudio) // 删除音频
}

// 用户相关路由
func setupUserRoutes(rg *gin.RouterGroup) {
	rg.PUT("/user", controllers.UpdateUser)                // 更新用户信息
	rg.PUT("/user/admin/:id", controllers.UpdateUserAdmin) // 更新用户权限
	rg.GET("/user", controllers.GetUserInfo)               // 获取当前登录的用户信息
	rg.DELETE("/user/:id", controllers.DeleteUser)         // 删除用户
}

// 设置相关路由
func setupSettingsRoutes(rg *gin.RouterGroup) {
	rg.PUT("/settings", controllers.UpdateSettings) // 更新系统设置
}

// Connect相关路由
func setupConnectRoutes(rg *gin.RouterGroup) {
	rg.POST("/addConnect", controllers.AddConnect)          // 添加 Connect
	rg.DELETE("/delConnect/:id", controllers.DeleteConnect) // 删除 Connect
}
