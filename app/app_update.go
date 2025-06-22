package app

import (
	"capital-tracker/lib/constant"

	tea "github.com/charmbracelet/bubbletea"
)

// update is where we handle input
func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch app.Screen {
		case constant.ModeMenu:
			cmd := app.Handler.Update_Menu(&app, msg)
			if cmd != nil {
				return app, cmd
			}
		case constant.ModeListTransaction:
			app.Handler.Update_ListTransaction(&app, msg)
		case constant.ModeCreateTransaction:
			app.Handler.Update_CreateTransaction(&app, msg)
		}
	}

	return app, nil
}
