package main

import (
	"capital-tracker/handler"
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
	cursor   int
	choices  []string
	selected string
	content  string
}

var h handler.Handler

func initialModel() model {
	return model{
		// choices to display
		choices: []string{"List Transactions", "Exit"},
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
		case "ctrl+c", "q":
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
				m.selected = "list_transaction"
				m.content = h.ListTransaction()
			case 1:
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

// view renders the UI
func (m model) View() string {
	s := "Use ↑ ↓ to move, Enter to select:\n\n"

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // current cursor
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	if m.selected != "" {
		s += "\n" + m.content
	}

	s += "\nPress q to quit."
	return s
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
