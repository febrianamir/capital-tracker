package repository

import "capital-tracker/model"

func (r *Repository) GetTransactions() []model.Transaction {
	var transactions []model.Transaction
	r.db.Model(&model.Transaction{}).Order("created_at DESC").Find(&transactions)
	return transactions
}
