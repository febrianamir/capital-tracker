package app

import (
	"capital-tracker/lib/constant"

	tea "github.com/charmbracelet/bubbletea"
)

// update is where we handle input
func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch m.Screen {
		case constant.ModeMenu:
			switch msg.String() {
			case "q":
				return m, tea.Quit

			case "up":
				if m.Cursor > 0 {
					m.Cursor--
				}

			case "down":
				if m.Cursor < len(m.Choices)-1 {
					m.Cursor++
				}

			case "enter":
				switch m.Cursor {
				case 0:
					m.Screen = constant.ModeListTransaction
				case 1:
					m.Screen = constant.ModeCreateTransaction
					m.CreateTransaction.FormValues = []string{}
					m.CreateTransaction.FormStep = 0
					m.CreateTransaction.CurrentInput = ""
				case 2:
					return m, tea.Quit
				}
			}
		case constant.ModeListTransaction:
			switch msg.Type {
			case tea.KeyCtrlC:
				m.Screen = constant.ModeMenu
			}
		case constant.ModeCreateTransaction:
			m.Handler.Update_CreateTransaction(&m, msg)
		}
	}

	return m, nil
}
