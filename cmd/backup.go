package cmd

import (
	"github.com/lin-snow/ech0/internal/cli"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "备份数据",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoBackup()
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
