package app

import (
	"capital-tracker/lib/constant"
)

// view renders the UI
func (app App) View() string {
	switch app.Screen {
	case constant.ModeMenu:
		return app.Handler.View_Menu(&app)
	case constant.ModeListTransaction:
		return app.Handler.View_ListTransaction(&app)
	case constant.ModeCreateTransaction:
		return app.Handler.View_CreateTransaction(&app)
	}

	return ""
}
