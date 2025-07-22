package cmd

import (
	"os"

	"github.com/lin-snow/ech0/internal/cli"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ech0",
	Short: "开源、自托管、专注思想流动的轻量级发布平台",
	Long:  `Ech0 是一款专为轻量级分享而设计的开源自托管平台，支持快速发布与分享你的想法、文字与链接。简单直观的操作界面，轻松管理你的内容，让分享变得更加自由，确保数据完全掌控，随时随地与世界连接。`,

	// 这个 Run 会在没有子命令时执行
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoServe()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看当前版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoVersion()
	},
}

func init() {
	// 解决Windows下使用 Cobra 触发 mousetrap 提示
    cobra.MousetrapHelpText = ""
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}


