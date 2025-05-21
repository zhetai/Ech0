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
	r.Use(static.Serve("/", static.LocalFile("./template", true)))
	r.Static("/api/images", "./data/images")
	// r.Static("/api/audios", "./data/audios")
	r.GET("/rss", controllers.GenerateRSS) // 生成 RSS 订阅链接

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

	// 需要鉴权的路由
	authRoutes := r.Group("/api")
	authRoutes.Use(middleware.JWTAuthMiddleware()) // 使用 JWT 鉴权中间件

	// 留言相关路由
	authRoutes.GET("/messages/:id", controllers.GetMessage)          // 获取留言详情
	authRoutes.POST("/messages/page", controllers.GetMessagesByPage) // 分页获取留言列表
	authRoutes.POST("/messages", controllers.PostMessage)            // 发布留言
	authRoutes.DELETE("/messages/:id", controllers.DeleteMessage)    // 删除留言
	authRoutes.POST("/images/upload", controllers.UploadImage)       // 上传图片
	authRoutes.DELETE("/images/delete", controllers.DeleteImage)     // 删除图片
	authRoutes.POST("/audios/upload", controllers.UploadAudio)       // 上传音频
	authRoutes.DELETE("/audios/delete", controllers.DeleteAudio)     // 删除音频

	// Todo 相关路由
	authRoutes.GET("/todo", controllers.GetTodos)          // 获取 Todo 列表
	authRoutes.POST("/todo", controllers.PostTodo)         // 发布 Todo
	authRoutes.PUT("/todo/:id", controllers.UpdateTodo)    // 更新 Todo
	authRoutes.DELETE("/todo/:id", controllers.DeleteTodo) // 删除 Todo

	// 用户相关路由
	authRoutes.PUT("/user", controllers.UpdateUser)                // 更新用户信息
	authRoutes.PUT("/user/admin/:id", controllers.UpdateUserAdmin) // 更新用户权限
	authRoutes.GET("/user", controllers.GetUserInfo)               // 获取当前登录的用户信息
	authRoutes.DELETE("/user/:id", controllers.DeleteUser)         // 删除用户

	// 设置相关路由
	authRoutes.PUT("/settings", controllers.UpdateSettings) // 更新系统设置

	// Connect相关路由
	authRoutes.POST("/addConnect", controllers.AddConnect)          // 添加 Connect
	authRoutes.DELETE("/delConnect/:id", controllers.DeleteConnect) // 删除 Connect

	// 其它路由

	// 由于Vue3 和SPA模式，所以处理匹配不到的路由(重定向到index.html)
	r.NoRoute(func(c *gin.Context) {
		c.File("./template/index.html")
	})

	return r
}
