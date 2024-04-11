package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	fakeEventData = []list.Item{
		EventEntry{
			title:       "Make moai app",
			description: "CODE CODE CODE",
			startTime:   time.Now().Add(time.Hour),
			endTime:     time.Now().Add(5 * time.Hour),
			isAllDay:    false,
		},
		EventEntry{
			title:       "Boo",
			description: "YA",
			startTime:   time.Now(),
			isAllDay:    true,
		},
	}

	allDayStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ECE6A")).
			Bold(true)
)

type EventEntry struct {
	title       string
	startTime   time.Time
	endTime     time.Time
	isAllDay    bool
	description string
}

func (entry EventEntry) Title() string {
	return entry.title
}
func (entry EventEntry) Description() string {
	text := strings.Builder{}
	if entry.isAllDay {
		text.WriteString(
			fmt.Sprintf("Date: %s ", entry.startTime.Format("02-01-06")),
		)
		//text.WriteString("Time: ")
		text.WriteString(allDayStyle.Render("ALL DAY\n"))
	}
	text.WriteString(
		fmt.Sprintf("Start: %s\n", entry.startTime.Format("02-01-06, 03:04 pm")),
	)
	text.WriteString(
		fmt.Sprintf("End:   %s\n", entry.startTime.Format("02-01-06, 03:04 pm")),
	)

	return text.String()
}
func (entry EventEntry) FilterValue() string {
	return entry.title + " " + entry.description
}
