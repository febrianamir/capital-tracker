package main

import (
	"capital-tracker/handler"
	"capital-tracker/lib/constant"
	dbmodel "capital-tracker/model"
	"capital-tracker/repository"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type model struct {
	cursor  int
	choices []string
	screen  constant.InputState
}

var h handler.Handler

func initialModel() model {
	return model{
		choices: []string{"List Transactions", "Exit"}, // menu list
		screen:  constant.ModeMenu,                     // default menu
	}
}

// messages are handled here
func (m model) Init() tea.Cmd {
	return nil
}

// update is where we handle input
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.screen = constant.ModeMenu

		case "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			switch m.cursor {
			case 0:
				m.screen = constant.ModeListTransaction
			case 1:
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

// view renders the UI
func (m model) View() string {
	switch m.screen {
	case constant.ModeMenu:
		s := "Use ↑ ↓ to move, Enter to select:\n\n"

		for i, choice := range m.choices {
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // current cursor
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += "\nPress q to quit."
		return s
	case constant.ModeListTransaction:
		return h.ListTransaction()
	}

	return ""
}

func init() {
	godotenv.Load()
}

func main() {
	db, err := gorm.Open(sqlite.Open("file:db/capital_tracker.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// migration
	db.AutoMigrate(&dbmodel.Transaction{})

	repo := repository.InitRepository(db)
	h = handler.InitHandler(&repo)

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
