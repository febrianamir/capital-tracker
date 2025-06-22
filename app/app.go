package app

import (
	"capital-tracker/lib/constant"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	Handler Handler

	Cursor  int
	Choices []string
	Screen  constant.InputState

	Menu              Menu
	ListTransaction   ListTransaction
	CreateTransaction CreateTransaction
}

type Menu struct {
	Content string
}

type ListTransaction struct {
	Cursor         int
	Choices        []string
	SelectedChoice string
}

type CreateTransaction struct {
	FormStep              int
	FormFields            []string
	FormFieldDescriptions []string
	FormValues            []string
	CurrentInput          string
}

func InitApp(handler Handler) App {
	tokens, _ := handler.repo.GetTransactionTokens()

	return App{
		Handler: handler,
		Choices: []string{ // menu list
			"List Transactions",
			"Create Transaction",
			"Exit",
		},
		Screen: constant.ModeMenu, // default menu

		ListTransaction: ListTransaction{
			Cursor:         0,
			Choices:        tokens,
			SelectedChoice: "",
		},
		CreateTransaction: CreateTransaction{
			FormFields:            []string{"Transaction Type", "Token", "Date", "Market Price", "Quantity", "Amount"},
			FormFieldDescriptions: []string{"Transaction Type (BUY/SELL)", "Token", "Date (DD/MM/YYYY HH:MM)", "Market Price", "Quantity", "Amount"},
			FormValues:            []string{},
		},
	}
}

// messages are handled here
func (m App) Init() tea.Cmd {
	return nil
}
