package app

import (
	"capital-tracker/lib/constant"
	"capital-tracker/response"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	Handler Handler

	Cursor  int
	Choices []string

	Screen  constant.InputState
	Spinner spinner.Model

	Menu              Menu
	ListTransaction   ListTransaction
	CreateTransaction CreateTransaction
}

type Menu struct {
	Content string
}

type ListTransaction struct {
	Cursor           int
	Choices          []string
	SelectedChoice   string
	IsLoading        bool
	CoinListResponse response.CoinList
	Error            error
}

type CreateTransaction struct {
	FormStep              int
	FormFields            []string
	FormFieldDescriptions []string
	FormValues            []string
	CurrentInput          string
}

type AppResponseMsg struct {
	CoinListResponse response.CoinList
	Error            error
}

func InitApp(handler Handler) App {
	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	tokens, _ := handler.repo.GetTransactionTokens()

	return App{
		Handler: handler,

		Cursor: 0,
		Choices: []string{ // menu list
			"List Transactions",
			"Create Transaction",
			"Exit",
		},

		Screen:  constant.ModeMenu, // default menu
		Spinner: s,

		ListTransaction: ListTransaction{
			Cursor:         0,
			Choices:        tokens,
			SelectedChoice: "",
		},
		CreateTransaction: CreateTransaction{
			FormStep: 0,
			FormFields: []string{
				"Transaction Type",
				"Token",
				"Date",
				"Market Price",
				"Quantity",
				"Amount",
			},
			FormFieldDescriptions: []string{
				"Transaction Type (BUY/SELL)",
				"Token",
				"Date (DD/MM/YYYY HH:MM)",
				"Market Price",
				"Quantity",
				"Amount",
			},
			FormValues:   []string{},
			CurrentInput: "",
		},
	}
}

// messages are handled here
func (app App) Init() tea.Cmd {
	return nil
}
