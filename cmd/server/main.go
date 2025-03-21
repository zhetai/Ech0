package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/config"
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/routers"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf(models.LoadConfigErrorMessage+": %v", err)
	}

	// 初始化数据库

	if err := database.InitDB(); err != nil {
		log.Fatalf(models.DatabaseInitErrorMessage+": %v", err)
	}

	// 设置Gin模式
	ginMode := config.Config.Server.Mode
	if ginMode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode) // 默认设置为Debug模式
	}

	// 设置路由
	r := routers.SetupRouter()

	// 启动服务器
	address := config.Config.Server.Host + ":" + config.Config.Server.Port
	if err := r.Run(address); err != nil {
		log.Fatalf(models.ServerLaunchErrorMessage+": %v", err)
	}
}
