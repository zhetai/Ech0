package cmd

import (
	"os"

	"github.com/lin-snow/ech0/internal/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ech0",
	Short: "开源、自托管、专注思想流动的轻量级发布平台",
	Long:  `开源、自托管、专注思想流动的轻量级发布平台`,

	// 这个 Run 会在没有子命令时执行
	Run: func(cmd *cobra.Command, args []string) {
		// 创建 Ech0 服务器
		s := server.New()

		// 初始化 Ech0
		s.Init()

		// 启动 Ech0
		s.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
