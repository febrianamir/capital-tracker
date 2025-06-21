package main

import (
	"capital-tracker/app"
	"capital-tracker/model"
	"capital-tracker/repository"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	godotenv.Load()
}

func main() {
	db, err := gorm.Open(sqlite.Open("file:db/capital_tracker.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// migration
	db.AutoMigrate(&model.Transaction{})

	repo := repository.InitRepository(db)
	h := app.InitHandler(&repo)
	app := app.InitApp(h)

	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
