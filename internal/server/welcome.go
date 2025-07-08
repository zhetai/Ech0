package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
)

var (
	isDarkBg = lipgloss.HasDarkBackground()

	// 信息样式（每行）
	infoStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.AdaptiveColor{
			Light: "236", Dark: "252",
		})

	// 高亮样式
	highlight = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{
			Light: "#0000FF", Dark: "#87CEFA",
		})

	// 外框
	boxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ff7675")).
			Padding(1, 4).
			Margin(1, 2)
)

// PrintGreetings 使用 lipgloss 输出欢迎信息
func PrintGreetings(port string) {
	// 渐变 Banner 渲染（每行变色）
	banner := gradientBanner(commonModel.GreetingBanner)

	// 构建正文内容
	content := lipgloss.JoinVertical(lipgloss.Left,
		infoStyle.Render("📦 Version: "+highlight.Render(commonModel.Version)),
		infoStyle.Render("🎈 Port: "+highlight.Render(port)),
		infoStyle.Render("🧙 Author: "+highlight.Render("L1nSn0w")),
		infoStyle.Render("👉 Website: "+highlight.Render("https://echo.soopy.cn/")),
		infoStyle.Render("👉 GitHub: "+highlight.Render("https://github.com/lin-snow/Ech0")),
	)

	full := lipgloss.JoinVertical(lipgloss.Left,
		banner,
		boxStyle.Render(content),
	)

	fmt.Fprintln(os.Stdout, full)
}

func gradientBanner(banner string) string {
	lines := strings.Split(banner, "\n")
	var rendered []string

	colors := []string{"#00BFFF", "#7B68EE", "#DA70D6", "#FF69B4", "#FF8C00", "#FFD700", "#00FA9A"}

	for i, line := range lines {
		color := lipgloss.Color(colors[i%len(colors)])
		style := lipgloss.NewStyle().Foreground(color)
		rendered = append(rendered, style.Render(line))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}
