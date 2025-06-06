package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/controllers"
)

func setupPublicRoutes(r *gin.Engine) {
	// 公共的路由
	publicRoutes := r.Group("/api")

	publicRoutes.POST("/login", controllers.Login)                  // 登录
	publicRoutes.POST("/register", controllers.Register)            // 注册
	publicRoutes.GET("/status", controllers.GetStatus)              // 获取用户信息
	publicRoutes.GET("/settings", controllers.GetSettings)          // 获取系统设置
	publicRoutes.GET("/heatmap", controllers.GetHeatMap)            // 获取热力图数据
	publicRoutes.GET("/allusers", controllers.GetAllUsers)          // 获取所有用户
	publicRoutes.GET("/connect", controllers.GetConnect)            // 获取Connect信息
	publicRoutes.GET("/connect/list", controllers.GetConnects)      // 获取添加的 Connect 列表
	publicRoutes.GET("/connects/info", controllers.GetConnectsInfo) // 获取 Connect信息列表
	publicRoutes.GET("/getmusic", controllers.GetPlayMusic)         // 获取音乐播放链接
	publicRoutes.GET("/playmusic", controllers.PlayMusic)           // 播放音乐

	// publicRoutes.GET("/messages", controllers.GetMessages) // 获取留言列表
}
