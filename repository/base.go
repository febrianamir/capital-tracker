package repository

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func InitRepository(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}
