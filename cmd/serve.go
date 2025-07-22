package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lin-snow/ech0/internal/cli"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 Web 服务",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoServe()

		// 阻塞主线程，直到接收到终止信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
