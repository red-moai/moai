package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/Genekkion/moai/components"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	COLUMN_WIDTH = 4
	DAYS_OF_WEEK = []table.Column{
		{
			Title: "Mon", Width: COLUMN_WIDTH,
		},
		{
			Title: "Tue", Width: COLUMN_WIDTH,
		},
		{
			Title: "Wed", Width: COLUMN_WIDTH,
		},
		{
			Title: "Thu", Width: COLUMN_WIDTH,
		},
		{
			Title: "Fri", Width: COLUMN_WIDTH,
		},
		{
			Title: "Sat", Width: COLUMN_WIDTH,
		},
		{
			Title: "Sun", Width: COLUMN_WIDTH,
		},
	}
)

type CalendarModel struct {
	currentTime *time.Time

	calendar    *table.Model
	targetMonth *time.Month
	targetYear  *int

	events *components.DefaultListModel
}

func InitCalendar(_ external.MoaiModel) tea.Model {
	currentTime := time.Now()
	targetMonth := currentTime.Month()
	targetYear := currentTime.Year()

	calendar := table.New(
		table.WithColumns(DAYS_OF_WEEK),
		table.WithFocused(false),
		table.WithHeight(6),
	)

	events := components.InitDefaultList(
		nil,
		"Events",
		30,
		15,
		nil,
	)

	model := CalendarModel{
		currentTime: &currentTime,
		calendar:    &calendar,
		targetMonth: &targetMonth,
		targetYear:  &targetYear,
		events:      &events,
	}

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FF0000")).
		BorderBottom(true).
		Bold(true)
	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#FFFFFF")).
		Bold(false)
	styles.Cell = styles.Cell.
		Bold(false)

	model.calendar.SetStyles(styles)

	return model
}

func (model CalendarModel) prettyDateTime() string {
	return fmt.Sprintf("%s, %s %d",
		model.currentTime.Weekday().String(),
		model.currentTime.Month().String(),
		//[:3],
		model.currentTime.Day(),
	)
}

func (model CalendarModel) Init() tea.Cmd {
	return nil
}

func (model CalendarModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	currentTime := time.Now()
	model.currentTime = &currentTime

	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "tab":
			if model.calendar.Focused() {
				model.calendar.Blur()
			} else {
				model.calendar.Focus()
			}
		}
	}

	var command tea.Cmd
	if model.calendar.Focused() {
		*model.calendar, command = model.calendar.Update(message)

	} else {
		*model.events, command = model.events.Update(message)
	}
	return model, command
}

func (model CalendarModel) View() string {
	text := strings.Builder{}
	text.WriteString(model.prettyDateTime() + "\n")

	rows := []table.Row{}
	numDays := model.getNumDays()
	x := 1
	row := table.Row{}
	for range model.getStartOffset() {
		row = append(row, "")
	}
	for len(row) < 7 {
		row = append(row, fmt.Sprintf("%d", x))
		x++
	}
	rows = append(rows, row)
	row = table.Row{}

	for x <= numDays {
		if len(row) < 7 {
			row = append(row, fmt.Sprintf("%d", x))
			x++
			continue
		}
		rows = append(rows, row)
		row = table.Row{}
	}

	if len(row) > 0 {
		for len(row) < 7 {
			row = append(row, "")
		}
		rows = append(rows, row)
	}
	model.calendar.SetRows(rows)

	text.WriteString(model.calendar.View() + "\n")
	text.WriteString(model.events.View() + "\n")
	return text.String()
}

// Returns the number of days before the 1st
// day in the month. Assumes calendar starts
// from Monday. Return value [0:6]
func (model *CalendarModel) getStartOffset() int {
	return (int(
		time.Date(
			*model.targetYear,
			*model.targetMonth,
			1, 0, 0, 0, 0,
			model.currentTime.Location(),
		).Weekday(),
	) + 6) % 7
}

func (model CalendarModel) getNumDays() int {
	month := *model.targetMonth
	year := *model.targetYear

	switch month {
	case time.February:
		if year%4 == 0 && (year%10 != 0 || year%400 == 0) {
			return 29
		}
		return 28
	case time.April, time.June, time.September, time.November:
		return 30
	}
	return 31
}
