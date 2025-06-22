package app

import (
	"capital-tracker/lib/constant"
	"fmt"
)

// view renders the UI
func (m App) View() string {
	switch m.Screen {
	case constant.ModeMenu:
		s := "Use ↑ ↓ to move, Enter to select:\n\n"

		for i, choice := range m.Choices {
			cursor := " " // no cursor
			if m.Cursor == i {
				cursor = ">" // current cursor
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		if m.Menu.Content != "" {
			s += "\n" + m.Menu.Content
		}

		s += "\nPress q to quit."
		return s
	case constant.ModeListTransaction:
		return m.Handler.View_ListTransaction(&m)
	case constant.ModeCreateTransaction:
		prompt := m.CreateTransaction.FormFieldDescriptions[m.CreateTransaction.FormStep]
		return fmt.Sprintf("Enter %s:\n%s\n\n(Press Enter to continue)", prompt, m.CreateTransaction.CurrentInput)
	}

	return ""
}
