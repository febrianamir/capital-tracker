package main

import (
	"capital-tracker/handler"
	"capital-tracker/lib/constant"
	"capital-tracker/model"
	"capital-tracker/repository"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type menu struct {
	content string
}

type createTransaction struct {
	formStep              int
	formFields            []string
	formFieldDescriptions []string
	formValues            []string
	currentInput          string
}

type app struct {
	cursor  int
	choices []string
	screen  constant.InputState

	menu              menu
	createTransaction createTransaction
}

var h handler.Handler
var repo repository.Repository

func initialApp() app {
	return app{
		choices: []string{ // menu list
			"List Transactions",
			"Create Transaction",
			"Exit",
		},
		screen: constant.ModeMenu, // default menu

		createTransaction: createTransaction{
			formFields:            []string{"Transaction Type", "Token", "Date", "Market Price", "Quantity", "Amount"},
			formFieldDescriptions: []string{"Transaction Type (BUY/SELL)", "Token", "Date (DD/MM/YYYY HH:MM)", "Market Price", "Quantity", "Amount"},
			formValues:            []string{},
		},
	}
}

// messages are handled here
func (m app) Init() tea.Cmd {
	return nil
}

// update is where we handle input
func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch m.screen {
		case constant.ModeMenu:
			switch msg.String() {
			case "q":
				return m, tea.Quit

			case "up":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}

			case "enter":
				switch m.cursor {
				case 0:
					m.screen = constant.ModeListTransaction
				case 1:
					m.screen = constant.ModeCreateTransaction
					m.createTransaction.formValues = []string{}
					m.createTransaction.formStep = 0
					m.createTransaction.currentInput = ""
				case 2:
					return m, tea.Quit
				}
			}
		case constant.ModeListTransaction:
			switch msg.Type {
			case tea.KeyCtrlC:
				m.screen = constant.ModeMenu
			}
		case constant.ModeCreateTransaction:
			switch msg.Type {
			case tea.KeyCtrlC:
				m.screen = constant.ModeMenu
			case tea.KeyEnter:
				// save current input
				m.createTransaction.formValues = append(m.createTransaction.formValues, m.createTransaction.currentInput)
				m.createTransaction.currentInput = ""
				m.createTransaction.formStep++

				if m.createTransaction.formStep >= len(m.createTransaction.formFields) {
					transaction := model.Transaction{}
					// validate & convert input data
					validationErrMsgList := ""
					for i, formValue := range m.createTransaction.formValues {
						// transaction type
						if i == 0 {
							msg := fmt.Sprintf("- %s is required, has to be BUY/SELL\n", m.createTransaction.formFields[i])
							validationErrMsg := ""

							if formValue == "" {
								validationErrMsg = msg
							}
							if validationErrMsg == "" && formValue != "BUY" && formValue != "SELL" {
								validationErrMsg = msg
							}

							validationErrMsgList += validationErrMsg
							transaction.TransactionType = formValue
						}

						// token
						if i == 1 {
							msg := fmt.Sprintf("- %s is required\n", m.createTransaction.formFields[i])
							validationErrMsg := ""

							if formValue == "" {
								validationErrMsg = msg
							}

							validationErrMsgList += validationErrMsg
							transaction.Token = formValue
						}

						// date
						if i == 2 {
							msg := fmt.Sprintf("- %s is required, and has to follow (DD/MM/YYYY HH:MM) format\n", m.createTransaction.formFields[i])
							validationErrMsg := ""
							layout := "02/01/2006 15:04"

							if formValue == "" {
								validationErrMsg = msg
							}
							_, err := time.Parse(layout, formValue)
							if validationErrMsg == "" && err != nil {
								validationErrMsg = msg
							}

							validationErrMsgList += validationErrMsg
							transaction.Date = formValue
						}

						// market price
						if i == 3 {
							msg := fmt.Sprintf("- %s is required, and has to be decimal number\n", m.createTransaction.formFields[i])
							validationErrMsg := ""

							if formValue == "" {
								validationErrMsg = msg
							}
							marketPrice, err := strconv.ParseFloat(formValue, 64)
							if validationErrMsg == "" && err != nil {
								validationErrMsg = msg
							}

							validationErrMsgList += validationErrMsg
							transaction.MarketPrice = marketPrice
						}

						// quantity
						if i == 4 {
							msg := fmt.Sprintf("- %s is required, and has to be decimal number\n", m.createTransaction.formFields[i])
							validationErrMsg := ""

							if formValue == "" {
								validationErrMsg = msg
							}
							quantity, err := strconv.ParseFloat(formValue, 64)
							if validationErrMsg == "" && err != nil {
								validationErrMsg = msg
							}

							validationErrMsgList += validationErrMsg
							transaction.Quantity = quantity
						}

						// amount
						if i == 5 {
							msg := fmt.Sprintf("- %s is required, and has to be decimal number\n", m.createTransaction.formFields[i])
							validationErrMsg := ""

							if formValue == "" {
								validationErrMsg = msg
							}
							amount, err := strconv.ParseFloat(formValue, 64)
							if validationErrMsg == "" && err != nil {
								validationErrMsg = msg
							}

							validationErrMsgList += validationErrMsg
							transaction.Amount = amount
						}
					}
					isValidationPass := validationErrMsgList == ""
					if !isValidationPass {
						m.menu.content = "❌ Transaction not saved!\n" + validationErrMsgList
						m.screen = constant.ModeMenu
					}

					// save transaction data to db
					if isValidationPass {
						_, err := repo.CreateTransaction(transaction)
						if err != nil {
							m.menu.content = fmt.Sprintf("❌ [ERROR] repository.create_transaction: %s", err.Error())
						} else {
							m.menu.content = "✅ Transaction saved!\n" + strings.Join(m.createTransaction.formValues, " | ")
						}

						m.screen = constant.ModeMenu
					}
				}
			case tea.KeyBackspace:
				if len(m.createTransaction.currentInput) > 0 {
					m.createTransaction.currentInput = m.createTransaction.currentInput[:len(m.createTransaction.currentInput)-1]
				}
			default:
				m.createTransaction.currentInput += msg.String()
			}
		}
	}

	return m, nil
}

// view renders the UI
func (m app) View() string {
	switch m.screen {
	case constant.ModeMenu:
		s := "Use ↑ ↓ to move, Enter to select:\n\n"

		for i, choice := range m.choices {
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // current cursor
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		if m.menu.content != "" {
			s += "\n" + m.menu.content
		}

		s += "\nPress q to quit."
		return s
	case constant.ModeListTransaction:
		return h.ListTransaction()
	case constant.ModeCreateTransaction:
		prompt := m.createTransaction.formFieldDescriptions[m.createTransaction.formStep]
		return fmt.Sprintf("Enter %s:\n%s\n\n(Press Enter to continue)", prompt, m.createTransaction.currentInput)
	}

	return ""
}

func init() {
	godotenv.Load()
}

func main() {
	db, err := gorm.Open(sqlite.Open("file:db/capital_tracker.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// migration
	db.AutoMigrate(&model.Transaction{})

	repo = repository.InitRepository(db)
	h = handler.InitHandler(&repo)

	p := tea.NewProgram(initialApp())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
