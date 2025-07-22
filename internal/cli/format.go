package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/lipgloss"
)

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

	// 外框
	boxStyle = lipgloss.NewStyle().
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#fb5151ff")).
			Padding(1, 1).
			Margin(1, 1)
)

func PrintCLIInfo(title, msg string) {
	// 使用 lipgloss 渲染 CLI 信息
	fmt.Fprintln(os.Stdout, infoStyle.Render(titleStyle.Render(title)+": "+highlight.Render(msg)))
}

func PrintCLIWithBox(items ...struct{ title, msg string }) {
	if len(items) == 0 {
		return
	}

	var content string
	for i, item := range items {
		line := infoStyle.Render(titleStyle.Render(item.title) + ": " + highlight.Render(item.msg))
		if i > 0 {
			content += "\n"
		}
		content += line
	}

	boxedContent := boxStyle.Render(content)
	fmt.Fprintln(os.Stdout, boxedContent)
}

func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls") // Windows 清屏命令
	} else {
		cmd = exec.Command("clear") // Linux/macOS 清屏命令
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
