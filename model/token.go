package model

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	Symbol      string
	ApiId       string
	DisplayName string
	IsWatchlist bool
}
