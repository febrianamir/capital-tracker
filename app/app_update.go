package app

import (
	"capital-tracker/lib/constant"

	"github.com/charmbracelet/bubbles/spinner"
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
				return app, tea.Batch(cmd)
			}

		case constant.ModeListWatchlist:
			cmds := app.Handler.Update_ListWatchlist(&app, msg)
			if cmds != nil {
				return app, tea.Batch(cmds...)
			}

		case constant.ModeListTransaction:
			cmds := app.Handler.Update_ListTransaction(&app, msg)
			if cmds != nil {
				return app, tea.Batch(cmds...)
			}

		case constant.ModeCreateTransaction:
			app.Handler.Update_CreateTransaction(&app, msg)
		}

	case AppResponseMsg:
		switch app.Screen {
		case constant.ModeListTransaction:
			app.ListTransaction.IsLoading = false
			app.ListTransaction.CoinListResponse = msg.CoinListResponse
			app.ListTransaction.Error = msg.Error

		case constant.ModeListWatchlist:
			app.ListWatchlist.IsLoading = false
			app.ListWatchlist.CoinListResponse = msg.CoinListResponse
			app.ListWatchlist.Error = msg.Error
		}

	case spinner.TickMsg:
		switch app.Screen {
		case constant.ModeListTransaction:
			if app.ListTransaction.IsLoading {
				var cmd tea.Cmd
				app.Spinner, cmd = app.Spinner.Update(msg)
				return app, tea.Batch(cmd)
			}

		case constant.ModeListWatchlist:
			if app.ListWatchlist.IsLoading {
				var cmd tea.Cmd
				app.Spinner, cmd = app.Spinner.Update(msg)
				return app, tea.Batch(cmd)
			}
		}
	}

	return app, nil
}
