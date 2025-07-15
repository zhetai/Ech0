package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
)

const (
	// GreetingBanner 是控制台横幅
	GreetingBanner = `
███████╗     ██████╗    ██╗  ██╗     ██████╗ 
██╔════╝    ██╔════╝    ██║  ██║    ██╔═████╗
█████╗      ██║         ███████║    ██║██╔██║
██╔══╝      ██║         ██╔══██║    ████╔╝██║
███████╗    ╚██████╗    ██║  ██║    ╚██████╔╝
╚══════╝     ╚═════╝    ╚═╝  ╚═╝     ╚═════╝ 
                                             
`
)

var (
	isDarkBg = lipgloss.HasDarkBackground()

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
			Light: "#0000FF", Dark: "#F6C177",
		})

	// 高亮样式
	highlight = lipgloss.NewStyle().
			Bold(false).
			Foreground(lipgloss.AdaptiveColor{
			Light: "#0000FF", Dark: "#87CEFA",
		})

	// 外框
	boxStyle = lipgloss.NewStyle().
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ff7675")).
			Padding(1, 4).
			Margin(1, 2)
)

// PrintGreetings 使用 lipgloss 输出欢迎信息
func PrintGreetings(port string) {
	// 渐变 Banner 渲染（每行变色）
	banner := gradientBanner(GreetingBanner)

	// 构建正文内容
	content := lipgloss.JoinVertical(lipgloss.Left,
		infoStyle.Render("📦 "+titleStyle.Render("Version")+": "+highlight.Render(commonModel.Version)),
		infoStyle.Render("🎈 "+titleStyle.Render("Port")+": "+highlight.Render(port)),
		infoStyle.Render("🧙 "+titleStyle.Render("Author")+": "+highlight.Render("L1nSn0w")),
		infoStyle.Render("👉 "+titleStyle.Render("Website")+": "+highlight.Render("https://echo.soopy.cn/")),
		infoStyle.Render("👉 "+titleStyle.Render("GitHub")+": "+highlight.Render("https://github.com/lin-snow/Ech0")),
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

	// Rose Pine Moon 渐变色调：由暗到亮，典雅梦幻
	colors := []string{
		"#232136", // Base
		"#393552", // Surface
		"#6E6A86", // Overlay
		"#EA9A97", // Love
		"#C4A7E7", // Iris
		"#9CCFD8", // Foam
		"#F6C177", // Gold
	}

	for i, line := range lines {
		color := lipgloss.Color(colors[i%len(colors)])
		style := lipgloss.NewStyle().Foreground(color)
		rendered = append(rendered, style.Render(line))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}
