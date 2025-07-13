package lib

import (
	"capital-tracker/lib/constant"
	"capital-tracker/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

type responseBody any

func DoRequest[T responseBody](method string, path string, queryParams map[string]string) (T, error) {
	var rspBody T

	queryParamsStr := buildQueryParams(queryParams)
	url := fmt.Sprintf("%s%s%s", os.Getenv("API_URL"), path, queryParamsStr)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return rspBody, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", os.Getenv("API_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return rspBody, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return rspBody, err
	}

	rspBody, err = parseResponse[T](body)
	return rspBody, err
}

func buildQueryParams(queryParams map[string]string) string {
	queryParamsStr := ""
	if len(queryParams) > 0 {
		queryParamsStr += "?"
		for key, value := range queryParams {
			queryParamsStr += fmt.Sprintf("%s=%s&", key, value)
		}
	}
	return queryParamsStr
}

func parseResponse[T responseBody](rspBodyByte []byte) (T, error) {
	var rspBody T
	err := json.Unmarshal(rspBodyByte, &rspBody)
	if err != nil {
		return rspBody, err
	}
	return rspBody, nil
}

func FormatPrice(printFormat string, price float64) string {
	// format to two decimal places
	priceStr := fmt.Sprintf(printFormat, price)

	// split integer and fractional parts
	parts := strings.Split(priceStr, ".")
	intPartStr := parts[0]
	var decimalPartStr string
	if len(parts) > 1 {
		decimalPartStr = parts[1]
	}

	// add comma formatting to the integer part
	intPart, _ := strconv.Atoi(intPartStr)
	intWithComma := humanize.Comma(int64(intPart))
	if strings.Contains(intWithComma, ".") {
		intWithComma = strings.Split(intWithComma, ".")[0]
	}

	// rejoin with decimal
	if decimalPartStr != "" {
		return intWithComma + "." + decimalPartStr
	}
	return intWithComma
}

func PrintLine(str string) string {
	return fmt.Sprintf("%s\n", str)
}

type CalculateHoldingStatRes struct {
	CostBasis          float64
	TotalQuantity      float64
	TotalCurrentAmount float64
	TotalPnlPercentage float64
	TotalPnlAmount     float64
}

func CalculateHoldingStats(transactions []model.Transaction, tokenPrice float64) (res CalculateHoldingStatRes) {
	totalBuyAmount := 0.0
	totalSellAmount := 0.0

	for _, transaction := range transactions {
		switch transaction.TransactionType {
		case constant.TransactionTypeBuy:
			totalBuyAmount += transaction.Amount
			res.TotalQuantity += transaction.Quantity
		case constant.TransactionTypeSell:
			totalSellAmount += transaction.Amount
			res.TotalQuantity -= transaction.Quantity
		}
	}

	if res.TotalQuantity <= 0 {
		res.TotalQuantity = 0
	}
	res.CostBasis += totalBuyAmount
	res.TotalCurrentAmount = res.TotalQuantity * tokenPrice
	res.TotalPnlAmount = (res.TotalCurrentAmount + totalSellAmount) - res.CostBasis
	res.TotalPnlPercentage = (res.TotalPnlAmount / res.CostBasis) * 100
	return res
}
