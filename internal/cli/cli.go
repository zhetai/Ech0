package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/lin-snow/ech0/internal/backup"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	"github.com/lin-snow/ech0/internal/server"
)

var s *server.Server // s 是全局的 Ech0 服务器实例

// DoServe 启动服务
func DoServe() {
	// 创建 Ech0 服务器
	s = server.New()
	// 初始化 Ech0
	s.Init()
	// 启动 Ech0
	s.Start()
}

// DoServeWithBlock 阻塞当前线程，直到服务器停止
func DoServeWithBlock() {
	// 创建 Ech0 服务器
	s = server.New()
	// 初始化 Ech0
	s.Init()
	// 启动 Ech0
	s.Start()

	// 阻塞主线程，直到接收到终止信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 创建 context，最大等待 5 秒优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		PrintCLIInfo("❌ 服务停止", "服务器强制关闭")
		os.Exit(1)
	}
	PrintCLIInfo("🎉 停止服务成功", "Ech0 服务器已停止")
}

// DoStopServe 停止服务
func DoStopServe() {
	if s == nil {
		PrintCLIInfo("⚠️ 停止服务", "Ech0 服务器未启动")
		return
	}

	// 创建 context，最大等待 5 秒优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		PrintCLIInfo("😭 停止服务失败", err.Error())
		return
	}

	s = nil // 清空全局服务器实例

	PrintCLIInfo("🎉 停止服务成功", "Ech0 服务器已停止")
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

// DoEch0Info() 打印 Ech0 信息
func DoEch0Info() {
	content := lipgloss.JoinVertical(lipgloss.Left,
		infoStyle.Render("📦 "+titleStyle.Render("Version")+": "+highlight.Render(commonModel.Version)),
		infoStyle.Render("🧙 "+titleStyle.Render("Author")+": "+highlight.Render("L1nSn0w")),
		infoStyle.Render("👉 "+titleStyle.Render("Website")+": "+highlight.Render("https://echo.soopy.cn/")),
		infoStyle.Render("👉 "+titleStyle.Render("GitHub")+": "+highlight.Render("https://github.com/lin-snow/Ech0")),
	)

	full := lipgloss.JoinVertical(lipgloss.Left,
		boxStyle.Render(content),
	)

	if _, err := fmt.Fprintln(os.Stdout, full); err != nil {
		fmt.Fprintf(os.Stderr, "failed to print ech0 info: %v\n", err)
	}
}

func DoHello() {
	ClearScreen()
	printCLIBanner()
}

// DoTui 执行 TUI
func DoTui() {
	ClearScreen()
	printCLIBanner()

	for {
		// 输出一行空行
		fmt.Println()

		var action string
		options := []huh.Option[string]{}

		if s == nil {
			options = append(options, huh.NewOption("🪅 启动 Web 服务", "serve"))
		} else {
			options = append(options, huh.NewOption("🛑 停止 Web 服务", "stopserve"))
		}

		options = append(options,
			huh.NewOption("🦖 查看信息", "info"),
			huh.NewOption("📦 执行备份", "backup"),
			huh.NewOption("💾 恢复备份", "restore"),
			huh.NewOption("📌 查看版本", "version"),
			huh.NewOption("❌ 退出", "exit"),
		)

		err := huh.NewSelect[string]().
			Title("欢迎使用 Ech0 TUI .").
			Options(options...).
			Value(&action).
			WithTheme(huh.ThemeCatppuccin()).
			Run()

		if err != nil {
			log.Fatal(err)
		}

		switch action {
		case "serve":
			ClearScreen()
			DoServe()
		case "stopserve":
			ClearScreen()
			DoStopServe()
		case "info":
			ClearScreen()
			DoEch0Info()
		case "backup":
			DoBackup()
		case "restore":
			// 如果服务器已经启动，则先停止服务器
			if s != nil {
				PrintCLIInfo("⚠️ 警告", "恢复数据前请先停止服务器")
			} else {
				// 获取备份文件路径
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
			}
		case "version":
			ClearScreen()
			DoVersion()
		case "exit":
			fmt.Println("👋 感谢使用 Ech0 TUI，期待下次再见")
			return
		}
	}
}

const (
	banner = `
    ______     __    ____ 
   / ____/____/ /_  / __ \
  / __/ / ___/ __ \/ / / /
 / /___/ /__/ / / / /_/ / 
/_____/\___/_/ /_/\____/  
`
)

func printCLIBanner() {
	lines := strings.Split(banner, "\n")
	var rendered []string

	colors := []string{
		"#FF7F7F", // 珊瑚红
		"#FFB347", // 桃橙色
		"#FFEB9C", // 金黄色
		"#B8E6B8", // 薄荷绿
		"#87CEEB", // 天空蓝
		"#DDA0DD", // 梅花紫
		"#F0E68C", // 卡其色
	}

	for i, line := range lines {
		color := lipgloss.Color(colors[i%len(colors)])
		style := lipgloss.NewStyle().Foreground(color)
		rendered = append(rendered, style.Render(line))
	}
	gradientBanner := lipgloss.JoinVertical(lipgloss.Left, rendered...)

	full := lipgloss.JoinVertical(lipgloss.Left,
		gradientBanner,
	)

	if _, err := fmt.Fprintln(os.Stdout, full); err != nil {
		fmt.Fprintf(os.Stderr, "failed to print banner: %v\n", err)
	}
}
