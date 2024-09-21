package main

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
	"log"
	"os"
	"os/exec"
)

type model struct {
	notifications []string
	cursor        int
	selected      map[int]struct{}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.notifications)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our notifications
	for i, choice := range m.notifications {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func InitialModel() model {
	notifications, err := dunstHistory()
	if err != nil {
		log.Fatalf("err initializing model %w", err)
	}
	choices := lo.Map(notifications, func(item Notification, index int) string {
		return item.Body
	})
	return model{
		// Our to-do list is a grocery list
		notifications: choices,

		// A map which indicates which notifications are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `notifications` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

type DunstHistory struct {
	Type string                `json:"type"`
	Data [][]DunstNotification `json:"data"`
}
type DunstNotification struct {
	Body struct {
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"body"`
	Message struct {
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"message"`
	Summary struct {
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"summary"`
	Appname struct {
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"appname"`
	Category struct {
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"category"`
	DefaultActionName struct {
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"default_action_name"`
	IconPath struct {
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"icon_path"`
	Id struct {
		Type string `json:"type"`
		Data int    `json:"data"`
	} `json:"id"`
	Timestamp struct {
		Type string `json:"type"`
		Data int64  `json:"data"`
	} `json:"timestamp"`
	Timeout struct {
		Type string `json:"type"`
		Data int    `json:"data"`
	} `json:"timeout"`
	Progress struct {
		Type string `json:"type"`
		Data int    `json:"data"`
	} `json:"progress"`
}

type Notification struct {
	App  string
	Body string
}

func dunstHistory() ([]Notification, error) {
	cmd := exec.Command("dunstctl", "history")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var dunstHistory DunstHistory
	err = json.Unmarshal(output, &dunstHistory)
	if err != nil {
		return nil, err
	}

	notfs := lo.Map(dunstHistory.Data, func(items []DunstNotification, index int) []Notification {
		return lo.Map(items, func(item DunstNotification, index int) Notification {
			return Notification{
				App:  item.Appname.Data,
				Body: item.Body.Data,
			}
		})
	})
	nots := lo.Flatten(notfs)
	return nots, nil
}

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Alas, there's been an error")
		os.Exit(1)
	}
}
