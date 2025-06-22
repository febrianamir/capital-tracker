package app

import (
	"capital-tracker/lib/constant"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (h *Handler) Update_Menu(app *App, msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q":
		return tea.Quit

	case "up":
		if app.Cursor > 0 {
			app.Cursor--
		}

	case "down":
		if app.Cursor < len(app.Choices)-1 {
			app.Cursor++
		}

	case "enter":
		switch app.Cursor {
		case 0:
			app.Screen = constant.ModeListTransaction
		case 1:
			app.Screen = constant.ModeCreateTransaction
			app.CreateTransaction.FormValues = []string{}
			app.CreateTransaction.FormStep = 0
			app.CreateTransaction.CurrentInput = ""
		case 2:
			return tea.Quit
		}
	}

	return nil
}

func (h *Handler) View_Menu(app *App) string {
	s := "Use ↑ ↓ to move, Enter to select:\n\n"

	for i, choice := range app.Choices {
		cursor := " " // no cursor
		if app.Cursor == i {
			cursor = ">" // current cursor
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	if app.Menu.Content != "" {
		s += "\n" + app.Menu.Content
	}

	s += "\nPress q to quit."
	return s
}
