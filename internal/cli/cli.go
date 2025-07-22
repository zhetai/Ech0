package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/lin-snow/ech0/internal/backup"
)

// DoBackup 执行备份
func DoBackup() {
	backupFilePath, _, err := backup.ExecuteBackup()
	if err != nil {
		// 处理错误
		PrintCLIInfo("😭 执行结果", "备份失败: "+err.Error())
		return
	}
	PrintCLIInfo("🎉 备份成功", backupFilePath)
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

var (
	// 信息样式（每行）
	infoStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.AdaptiveColor{
			Light: "236", Dark: "252",
		})

	// 标题样式
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{
			Light: "#4338ca", Dark: "#FF7F7F",
		})

	// 高亮样式
	highlight = lipgloss.NewStyle().
			Bold(false).
			Italic(true).
			Foreground(lipgloss.AdaptiveColor{
			Light: "#7c3aed", Dark: "#53b7f5ff",
		})
)

func PrintCLIInfo(title, msg string) {
	// 使用 lipgloss 渲染 CLI 信息
	fmt.Fprintln(os.Stdout, infoStyle.Render(titleStyle.Render(title)+": "+highlight.Render(msg)))
}