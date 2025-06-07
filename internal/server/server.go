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

type Server struct {
	GinEngine *gin.Engine
}

func New() *Server {
	engine := gin.Default()
	return &Server{
		GinEngine: engine,
	}
}

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

func printGreetings(port string) {
	fmt.Print(commonModel.GreetingBanner)
	fmt.Printf("---\n📦 Version: %s\n", commonModel.Version)
	fmt.Printf("🎈 Port: %s\n", port)
	fmt.Printf("🧙 Author: L1nSn0w\n")
	fmt.Printf("👉 Website: https://echo.soopy.cn/\n")
	fmt.Printf("👉 GitHub: https://github.com/lin-snow/Ech0\n---\n")
}
