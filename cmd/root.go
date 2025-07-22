package cmd

import (
	"os"

	"github.com/lin-snow/ech0/internal/cli"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ech0",
	Short: "开源、自托管、专注思想流动的轻量级发布平台",
	Long:  `开源、自托管、专注思想流动的轻量级发布平台`,

	// 这个 Run 会在没有子命令时执行
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoTui()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看当前版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoVersion()
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "查看当前信息",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoEch0Info()
	},
}

func init() {
	// 解决Windows下使用 Cobra 触发 mousetrap 提示
	cobra.MousetrapHelpText = ""
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(infoCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
