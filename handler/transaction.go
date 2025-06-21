package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

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

type Quotes struct {
	Data map[string]QuoteData `json:"data"`
}

type QuoteData struct {
	Quote struct {
		USD struct {
			Price float64 `json:"price"`
		} `json:"USD"`
	} `json:"quote"`
}

func (h *Handler) ListTransaction() string {
	var builder strings.Builder

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	id := "1"
	q.Add("id", id)
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("CMC_API_KEY"))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	respBody, _ := io.ReadAll(resp.Body)

	var quotes Quotes
	if err := json.Unmarshal(respBody, &quotes); err != nil {
		log.Fatal("Failed to unmarshal JSON:", err)
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

	tokenPrice := quotes.Data[id].Quote.USD.Price
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
