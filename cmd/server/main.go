package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/config"
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/routers"
)

func printGreetings(port string) {
	fmt.Printf("---\nGin Server Starting\nport: %s\n---\n", port)
	fmt.Print(models.GreetingBanner)
	fmt.Printf("Server has started on port %s\n", port)
	fmt.Printf("---\n📦 Version: %s\n", models.Version)
	fmt.Printf("🧙 Author: L1nSn0w\n")
	fmt.Printf("👉 Website: https://echo.soopy.cn/\n")
	fmt.Printf("👉 GitHub: https://github.com/lin-snow/Ech0\n---\n")
}

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
	r := gin.Default()
	routers.SetupRouter(r)

	// 启动服务器
	address := config.Config.Server.Host + ":" + config.Config.Server.Port
	printGreetings(config.Config.Server.Port)
	if err := r.Run(address); err != nil {
		log.Fatalf(models.ServerLaunchErrorMessage+": %v", err)
	}
}
