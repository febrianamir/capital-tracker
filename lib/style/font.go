package style

import "github.com/charmbracelet/lipgloss"

func FontBold(str string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Render(str)
}

func FontItalic(str string) string {
	return lipgloss.NewStyle().
		Italic(true).
		Render(str)
}
