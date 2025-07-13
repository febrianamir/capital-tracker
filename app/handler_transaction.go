package app

import (
	"capital-tracker/lib"
	"capital-tracker/lib/constant"
	"capital-tracker/lib/style"
	"capital-tracker/model"
	"capital-tracker/param"
	"capital-tracker/response"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (h *Handler) Update_ListTransaction(app *App, msg tea.KeyMsg) (cmds []tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		if app.ListTransaction.SelectedChoice != "" {
			app.ListTransaction.SelectedChoice = ""
			return nil
		}
		if app.ListTransaction.SelectedChoice == "" {
			app.Screen = constant.ModeMenu
		}

	case tea.KeyUp:
		if app.ListTransaction.Cursor > 0 {
			app.ListTransaction.Cursor--
		}

	case tea.KeyDown:
		if app.ListTransaction.Cursor < len(app.ListTransaction.Choices)-1 {
			app.ListTransaction.Cursor++
		}

	case tea.KeyEnter:
		app.ListTransaction.SelectedChoice = app.ListTransaction.Choices[app.ListTransaction.Cursor]
		app.ListTransaction.IsLoading = true

		tokens := map[string]string{
			"BTC":  "bitcoin",
			"HYPE": "hyperliquid",
		}
		queryParams := map[string]string{
			"vs_currency": "usd",
			"ids":         tokens[app.ListTransaction.SelectedChoice],
			"precision":   "2",
		}

		cmds = append(cmds, app.Spinner.Tick, func() tea.Msg {
			coins, err := lib.DoRequest[response.CoinList](http.MethodGet, "/coins/markets", queryParams)
			if err != nil {
				return AppResponseMsg{Error: err}
			}

			return AppResponseMsg{CoinListResponse: coins}
		})
	}

	return cmds
}

func (h *Handler) View_ListTransaction(app *App) string {
	if app.ListTransaction.SelectedChoice == "" {
		s := "[List Transaction] Use ↑ ↓ to move, Enter to select:\n\n"
		for i, choice := range app.ListTransaction.Choices {
			cursor := " " // no cursor
			if app.ListTransaction.Cursor == i {
				cursor = ">" // current cursor
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		return s
	}

	if app.ListTransaction.IsLoading {
		return fmt.Sprintf("\n  %s Loading...\n\n", app.Spinner.View())
	}

	var builder strings.Builder

	token := app.ListTransaction.SelectedChoice
	transactions, err := h.repo.GetTransactions(param.GetTransactions{
		Token: token,
	})
	if err != nil {
		return fmt.Sprintf("[ERROR] repository.get_transactions: %s", err.Error())
	}

	// each item in renderedContents represent one line
	renderedContents := []string{}

	renderedContents = append(renderedContents, "-------------------------------------------------------------------------------------------")
	renderedContents = append(renderedContents, style.FontBold(style.ColorCyan(app.ListTransaction.SelectedChoice)))

	tokenPrice := app.ListTransaction.CoinListResponse[0].CurrentPrice
	tokenPriceFormatted := lib.FormatPrice("%g", tokenPrice)
	if tokenPrice > 100000 {
		tokenPriceFormatted = lib.FormatPrice("%.2f", tokenPrice)
	}
	renderedContents = append(renderedContents, fmt.Sprintf("Current Price: $%s", tokenPriceFormatted))

	holdingStat := lib.CalculateHoldingStats(transactions, tokenPrice)

	renderedContents = append(renderedContents, fmt.Sprintf("Cost Basis: $%s", lib.FormatPrice("%.2f", holdingStat.CostBasis)))
	renderedContents = append(renderedContents, fmt.Sprintf("Current Amount: $%s (%.8f %s)", lib.FormatPrice("%.2f", holdingStat.TotalCurrentAmount), holdingStat.TotalQuantity, app.ListTransaction.SelectedChoice))

	totalPnl := fmt.Sprintf("$%s (%.2f%%)", lib.FormatPrice("%.2f", holdingStat.TotalPnlAmount), holdingStat.TotalPnlPercentage)
	if holdingStat.TotalPnlPercentage > 0 {
		totalPnl = style.ColorGreen(totalPnl)
	}
	if holdingStat.TotalPnlPercentage < 0 {
		totalPnl = style.ColorRed(totalPnl)
	}

	renderedContents = append(renderedContents, fmt.Sprintf("PnL: %s", totalPnl))
	renderedContents = append(renderedContents, "-------------------------------------------------------------------------------------------")
	renderedContents = append(renderedContents, "")

	headerFormat := "%-7s %-20s %-10s %-16s %-15s %-11s"
	header := fmt.Sprintf(headerFormat, "", "Datetime", "Token", "Market Price", "Quantity", "Amount")
	renderedContents = append(renderedContents, style.FontBold(style.ColorCyan(header)))

	dataFormat := "%-7s %-20s %-10s $%-15s %-15g $%-10s"

	for _, transaction := range transactions {
		transactionType := style.FontItalic(fmt.Sprintf("%-7s", fmt.Sprintf("(%s)", transaction.TransactionType)))
		if transaction.TransactionType == constant.TransactionTypeBuy {
			transactionType = style.ColorGreen(transactionType)
		}
		if transaction.TransactionType == constant.TransactionTypeSell {
			transactionType = style.ColorRed(transactionType)
		}
		renderedContents = append(renderedContents, fmt.Sprintf(dataFormat, transactionType, transaction.Date, transaction.Token, lib.FormatPrice("%g", transaction.MarketPrice), transaction.Quantity, lib.FormatPrice("%.2f", transaction.Amount)))
	}

	for _, renderedContent := range renderedContents {
		builder.WriteString(lib.PrintLine(renderedContent))
	}

	return builder.String()
}

func (h *Handler) Update_CreateTransaction(app *App, msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyCtrlC:
		app.Screen = constant.ModeMenu

	case tea.KeyEnter:
		// save current input
		app.CreateTransaction.FormValues = append(app.CreateTransaction.FormValues, app.CreateTransaction.CurrentInput)
		app.CreateTransaction.CurrentInput = ""
		app.CreateTransaction.FormStep++

		if app.CreateTransaction.FormStep >= len(app.CreateTransaction.FormFields) {
			transaction := model.Transaction{}
			// validate & convert input data
			validationErrMsgList := ""
			for i, formValue := range app.CreateTransaction.FormValues {
				// transaction type
				if i == 0 {
					msg := fmt.Sprintf("- %s is required, has to be BUY/SELL\n", app.CreateTransaction.FormFields[i])
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
					msg := fmt.Sprintf("- %s is required\n", app.CreateTransaction.FormFields[i])
					validationErrMsg := ""

					if formValue == "" {
						validationErrMsg = msg
					}

					validationErrMsgList += validationErrMsg
					transaction.Token = formValue
				}

				// date
				if i == 2 {
					msg := fmt.Sprintf("- %s is required, and has to follow (DD/MM/YYYY HH:MM) format\n", app.CreateTransaction.FormFields[i])
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
					msg := fmt.Sprintf("- %s is required, and has to be decimal number\n", app.CreateTransaction.FormFields[i])
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
					msg := fmt.Sprintf("- %s is required, and has to be decimal number\n", app.CreateTransaction.FormFields[i])
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
					msg := fmt.Sprintf("- %s is required, and has to be decimal number\n", app.CreateTransaction.FormFields[i])
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
				app.Menu.Content = "❌ Transaction not saved!\n" + validationErrMsgList
				app.Screen = constant.ModeMenu
			}

			// save transaction data to db
			if isValidationPass {
				_, err := h.repo.CreateTransaction(transaction)
				if err != nil {
					app.Menu.Content = fmt.Sprintf("❌ [ERROR] repository.create_transaction: %s", err.Error())
				} else {
					app.Menu.Content = "✅ Transaction saved!\n" + strings.Join(app.CreateTransaction.FormValues, " | ")
				}

				app.Screen = constant.ModeMenu
			}
		}

	case tea.KeyBackspace:
		if len(app.CreateTransaction.CurrentInput) > 0 {
			app.CreateTransaction.CurrentInput = app.CreateTransaction.CurrentInput[:len(app.CreateTransaction.CurrentInput)-1]
		}

	default:
		app.CreateTransaction.CurrentInput += msg.String()
	}
}

func (h *Handler) View_CreateTransaction(app *App) string {
	prompt := app.CreateTransaction.FormFieldDescriptions[app.CreateTransaction.FormStep]
	return fmt.Sprintf("Enter %s:\n%s\n\n(Press Enter to continue)", prompt, app.CreateTransaction.CurrentInput)
}
