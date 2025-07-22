package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/lin-snow/ech0/internal/backup"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	"github.com/lin-snow/ech0/internal/server"
)

// DoServe 启动服务
func DoServe() {
	// 创建 Ech0 服务器
	s := server.New()

	// 初始化 Ech0
	s.Init()

	// 启动 Ech0
	s.Start()
}

// DoBackup 执行备份
func DoBackup() {
	_, backupFileName, err := backup.ExecuteBackup()
	if err != nil {
		// 处理错误
		PrintCLIInfo("😭 执行结果", "备份失败: "+err.Error())
		return
	}

	// 获取PWD环境变量
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "backup", backupFileName)

	PrintCLIInfo("🎉 备份成功", fullPath)
}

// DoRestore 执行恢复
func DoRestore(backupFilePath string) {
	err := backup.ExecuteRestore(backupFilePath)
	if err != nil {
		// 处理错误
		PrintCLIInfo("😭 执行结果", "恢复失败: "+err.Error())
		return
	}
	PrintCLIInfo("🎉 恢复成功", "已从备份文件 "+backupFilePath+" 中恢复数据")
}

// DoVersion 打印版本信息
func DoVersion() {
	PrintCLIWithBox(struct{ title, msg string }{
		title: "📦 当前版本",
		msg:   "v" + commonModel.Version,
	})
}

// DoTui 执行 TUI
func DoTui() {
	for {
		var action string
		err := huh.NewSelect[string]().
			Title("🌟 Ech0 CLI 操作菜单").
			Options(
				huh.NewOption("🪅 启动 Web 服务", "serve"),
				huh.NewOption("📦 执行备份", "backup"),
				huh.NewOption("💾 恢复数据", "restore"),
				huh.NewOption("📌 查看版本信息", "version"),
				huh.NewOption("❌ 退出", "exit"),
			).
			Value(&action).
			Run()

		if err != nil {
			log.Fatal(err)
		}

		switch action {
		case "serve":
			DoServe()
		case "backup":
			DoBackup()
		case "restore":
			var path string
			huh.NewInput().
				Title("请输入备份文件路径").
				Value(&path).
				Run()
			path = strings.TrimSpace(path)
			if path != "" {
				DoRestore(path)
			} else {
				PrintCLIInfo("⚠️ 跳过", "未输入备份路径")
			}
		case "version":
			DoVersion()
		case "exit":
			fmt.Println("👋 已退出")
			return
		}
	}
}