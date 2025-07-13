package repository

import (
	"capital-tracker/lib"
	"capital-tracker/model"
	"capital-tracker/param"
	"context"
	"os"
)

func (r *Repository) GetTokens(param param.GetTokens) ([]model.Token, error) {
	var tokens []model.Token
	db, _ := lib.WithSQLLogger(context.Background(), r.db, lib.WithEnv(os.Getenv("ENV")))

	query := db.Model(&model.Token{})
	if param.IsWatchlist != nil {
		query.Where("is_watchlist = ?", param.IsWatchlist)
	}
	query = query.Order("created_at ASC")

	err := query.Find(&tokens).Error
	return tokens, err
}
