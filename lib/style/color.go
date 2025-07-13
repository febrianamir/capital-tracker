package style

import "github.com/charmbracelet/lipgloss"

func ColorGreen(str string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A3BE8C")).
		Render(str)
}

func ColorRed(str string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BF616A")).
		Render(str)
}

func ColorBlue(str string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#81A1C1")).
		Render(str)
}

func ColorCyan(str string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88C0D0")).
		Render(str)
}
