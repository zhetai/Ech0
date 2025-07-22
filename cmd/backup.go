package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "备份数据",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("正在备份数据...")
		// 调用 internal/backup 模块中的逻辑
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
