package repository

import (
	"capital-tracker/model"
	"capital-tracker/param"
)

func (r *Repository) GetTransactions(param param.GetTransactions) ([]model.Transaction, error) {
	var transactions []model.Transaction

	query := r.db.Model(&model.Transaction{})
	if param.Token != "" {
		query.Where("token = ?", param.Token)
	}
	query = query.Order("created_at DESC")

	err := query.Find(&transactions).Error
	return transactions, err
}

func (r *Repository) GetTransactionTokens() ([]string, error) {
	var tokens []string

	err := r.db.
		Model(&model.Transaction{}).
		Distinct("token").
		Pluck("token", &tokens).
		Error

	return tokens, err
}

func (r *Repository) CreateTransaction(data model.Transaction) (model.Transaction, error) {
	err := r.db.Model(&model.Transaction{}).Create(&data).Error
	return data, err
}
