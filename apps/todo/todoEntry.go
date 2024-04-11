package todo

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
)

var (
	fakeTodoData = []list.Item{
		TodoEntry{
			title:       "Play Content Warning 🎬",
			description: "With the fellas",
			deadline:    time.Now().Add(time.Hour),
			completed:   false,
			completedAt: nil,
		},
		TodoEntry{
			title:       "Crochet for someone 🧶",
			description: "Make cool stuff",
			deadline:    time.Now().Add(time.Hour),
			completed:   false,
			completedAt: nil,
		},
	}
)

type TodoEntry struct {
	title       string
	description string
	deadline    time.Time
	completed   bool
	completedAt *time.Time
}

func (entry TodoEntry) Title() string {
	return entry.title
}
func (entry TodoEntry) Description() string {
	return entry.description
}
func (entry TodoEntry) FilterValue() string {
	return entry.title
}
