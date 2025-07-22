package cmd

import (
	"github.com/lin-snow/ech0/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 Web 服务",
	Run: func(cmd *cobra.Command, args []string) {
		// 创建 Ech0 服务器
		s := server.New()

		// 初始化 Ech0
		s.Init()

		// 启动 Ech0
		s.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}