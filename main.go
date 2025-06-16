package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor   int
	choices  []string
	selected string
}

func initialModel() model {
	return model{
		// choices to display
		choices: []string{"Transactions", "Exit"},
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
				m.selected = "List Transactions"
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
		s += fmt.Sprintf("\n%s\n", m.selected)
	}

	s += "\nPress q to quit."
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
