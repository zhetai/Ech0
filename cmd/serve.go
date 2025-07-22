package cmd

import (
	"github.com/lin-snow/ech0/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 Web 服务",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.New()
		s.Init()
		s.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}