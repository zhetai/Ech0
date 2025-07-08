package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/config"
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/di"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	"github.com/lin-snow/ech0/internal/router"
	errUtil "github.com/lin-snow/ech0/internal/util/err"
	logUtil "github.com/lin-snow/ech0/internal/util/log"
)

// Server 服务器结构体，包含Gin引擎
type Server struct {
	GinEngine *gin.Engine
}

// New 创建一个新的服务器实例
func New() *Server {
	return &Server{}
}

// Init 初始化服务器
func (s *Server) Init() {
	// Logger
	logUtil.InitLogger()

	// Config
	config.LoadAppConfig()

	// Mode
	if config.Config.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Gin Engine
	s.GinEngine = gin.New()

	// Database
	database.InitDatabase()

	// Handlers
	handlers, err := di.BuildHandlers(database.DB)
	if err != nil {
		errUtil.HandlePanicError(&commonModel.ServerError{
			Msg: commonModel.INIT_HANDLERS_PANIC,
			Err: err,
		})
	}

	// Router
	router.SetupRouter(s.GinEngine, handlers)
}

// Start 启动服务器
func (s *Server) Start() {
	port := config.Config.Server.Port
	printGreetings(port)
	if err := s.GinEngine.Run(":" + port); err != nil {
		errUtil.HandlePanicError(&commonModel.ServerError{
			Msg: commonModel.GIN_RUN_FAILED,
			Err: err,
		})
	}
}

// printGreetings 打印欢迎信息
func printGreetings(port string) {
	fmt.Print(commonModel.GreetingBanner)
	fmt.Printf("---\n📦 Version: %s\n", commonModel.Version)
	fmt.Printf("🎈 Port: %s\n", port)
	fmt.Printf("🧙 Author: L1nSn0w\n")
	fmt.Printf("👉 Website: https://echo.soopy.cn/\n")
	fmt.Printf("👉 GitHub: https://github.com/lin-snow/Ech0\n---\n")
}
