package app

import (
	"capital-tracker/lib"
	"capital-tracker/lib/constant"
	"capital-tracker/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

func formatPrice(printFormat string, price float64) string {
	// format to two decimal places
	priceStr := fmt.Sprintf(printFormat, price)

	// split integer and fractional parts
	parts := strings.Split(priceStr, ".")
	intPartStr := parts[0]
	decimalPartStr := parts[1]

	// add comma formatting to the integer part
	intPart, _ := strconv.Atoi(intPartStr)
	intWithComma := humanize.Comma(int64(intPart))
	if strings.Contains(intWithComma, ".") {
		intWithComma = strings.Split(intWithComma, ".")[0]
	}

	// rejoin with decimal
	return intWithComma + "." + decimalPartStr
}

func printLine(str string) string {
	return fmt.Sprintf("%s\n", str)
}

type CoinList []Coin

type Coin struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"current_price"`
}

func (h *Handler) ListTransaction() string {
	var builder strings.Builder

	coins, err := lib.DoRequest[CoinList](http.MethodGet, "/coins/markets", map[string]string{
		"vs_currency": "usd",
		"ids":         "bitcoin",
		"precision":   "2",
	})
	if err != nil {
		return fmt.Sprintln(err.Error())
	}

	colorCyan := color.New(color.FgCyan).SprintFunc()
	colorRed := color.New(color.FgRed).SprintFunc()
	colorGreen := color.New(color.FgGreen).SprintFunc()
	styleBold := color.New(color.Bold).SprintFunc()
	styleItalic := color.New(color.Italic).SprintFunc()

	transactions, err := h.repo.GetTransactions()
	if err != nil {
		return fmt.Sprintf("[ERROR] repository.get_transactions: %s", err.Error())
	}

	// each item in renderedContents represent one line
	renderedContents := []string{}

	renderedContents = append(renderedContents, "-------------------------------------------------------------------------------------------")
	renderedContents = append(renderedContents, styleBold(colorCyan("BTC")))

	tokenPrice := coins[0].CurrentPrice
	tokenPriceFormatted := formatPrice("%g", tokenPrice)
	if tokenPrice > 100000 {
		tokenPriceFormatted = formatPrice("%.2f", tokenPrice)
	}
	renderedContents = append(renderedContents, fmt.Sprintf("Current Price: $%s", tokenPriceFormatted))

	costBasis := 0.0
	totalQuantity := 0.0
	for _, transaction := range transactions {
		costBasis += transaction.Amount
		totalQuantity += transaction.Quantity
	}
	totalCurrentAmount := totalQuantity * tokenPrice

	renderedContents = append(renderedContents, fmt.Sprintf("Cost Basis: $%s", formatPrice("%.2f", costBasis)))
	renderedContents = append(renderedContents, fmt.Sprintf("Current Amount: $%s (%.8f BTC)", formatPrice("%.2f", totalCurrentAmount), totalQuantity))

	totalPnlPercentage := ((totalCurrentAmount - costBasis) / costBasis) * 100
	totalPnl := fmt.Sprintf("%.2f%%", totalPnlPercentage)
	if totalPnlPercentage > 0 {
		totalPnl = colorGreen(totalPnl)
	}
	if totalPnlPercentage < 0 {
		totalPnl = colorRed(totalPnl)
	}

	renderedContents = append(renderedContents, fmt.Sprintf("PnL: %s", totalPnl))
	renderedContents = append(renderedContents, "-------------------------------------------------------------------------------------------")
	renderedContents = append(renderedContents, "")

	headerFormat := "%-7s %-20s %-10s %-16s %-15s %-11s"
	header := fmt.Sprintf(headerFormat, "", "Datetime", "Token", "Market Price", "Quantity", "Amount")
	renderedContents = append(renderedContents, styleBold(colorCyan(header)))

	dataFormat := "%-7s %-20s %-10s $%-15s %-15g $%-10s"

	for _, transaction := range transactions {
		transactionType := styleItalic(fmt.Sprintf("%-7s", fmt.Sprintf("(%s)", transaction.TransactionType)))
		if transaction.TransactionType == "BUY" {
			transactionType = colorGreen(transactionType)
		}
		if transaction.TransactionType == "SELL" {
			transactionType = colorRed(transactionType)
		}
		renderedContents = append(renderedContents, fmt.Sprintf(dataFormat, transactionType, transaction.Date, transaction.Token, formatPrice("%g", transaction.MarketPrice), transaction.Quantity, formatPrice("%.2f", transaction.Amount)))
	}

	for _, renderedContent := range renderedContents {
		builder.WriteString(printLine(renderedContent))
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
