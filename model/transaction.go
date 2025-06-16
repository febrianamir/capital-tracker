package model

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	TransactionType string
	Date            string
	Token           string
	MarketPrice     float64
	Quantity        float64
	Amount          float64
}
