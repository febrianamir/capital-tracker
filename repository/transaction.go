package repository

import "capital-tracker/model"

func (r *Repository) GetTransactions() ([]model.Transaction, error) {
	var transactions []model.Transaction

	query := r.db.Model(&model.Transaction{})
	query = query.Order("created_at DESC")

	err := query.Find(&transactions).Error
	return transactions, err
}

func (r *Repository) CreateTransaction(data model.Transaction) (model.Transaction, error) {
	err := r.db.Model(&model.Transaction{}).Create(&data).Error
	return data, err
}
