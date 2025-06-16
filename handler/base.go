package handler

import (
	"capital-tracker/repository"
)

type Handler struct {
	repo *repository.Repository
}

func InitHandler(repo *repository.Repository) Handler {
	return Handler{
		repo: repo,
	}
}
