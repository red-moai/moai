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

	styles      = table.DefaultStyles()
	headerStyle = styles.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#FF0000")).
			BorderBottom(true).
			Bold(true)

	defaultSelectedStyle = table.DefaultStyles().Selected.
				Foreground(lipgloss.Color("#FAF9F6")).
				Bold(false)

	activeSelectedStyle = table.DefaultStyles().Selected.
				Foreground(lipgloss.Color("#1A1B26")).
				Background(lipgloss.Color("#F7768E"))

	cellStyle = styles.Cell.
			Align(lipgloss.Center).
			Bold(false)

	todayStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9AA5CE")).
			Underline(true).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B4F9F8")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Margin(1, 0).
			Bold(true)

	modelStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1)
)

type CalendarModel struct {
	currentTime *time.Time

	focusedIndex *int
	calendar     *[]table.Model
	targetMonth  *time.Month
	targetYear   *int
	events       *components.DefaultListModel

	MainModel external.MoaiModel
}

func getDefaultStyle() table.Styles {
	styles.Selected = defaultSelectedStyle
	return styles
}

func getActiveStyle() table.Styles {
	styles.Selected = activeSelectedStyle
	return styles
}

func InitCalendar(mainModel external.MoaiModel) external.MoaiApp{
	currentTime := time.Now()
	targetMonth := currentTime.Month()
	targetYear := currentTime.Year()

	styles.Header = headerStyle
	styles.Cell = cellStyle

	defaultStyle := getDefaultStyle()

	calendar := make([]table.Model, 7)
	for i := range 7 {
		calendar[i] = table.New(
			table.WithColumns(
				[]table.Column{
					DAYS_OF_WEEK[i],
				},
			),
			table.WithFocused(false),
			table.WithHeight(6),
			table.WithStyles(defaultStyle),
		)

	}

	events := components.InitDefaultList(
		fakeEventData,
		"Events",
		mainModel,
		nil,
		nil,
	)

	focusedIndex := 0
	model := CalendarModel{
		currentTime:  &currentTime,
		calendar:     &calendar,
		targetMonth:  &targetMonth,
		targetYear:   &targetYear,
		events:       &events,
		focusedIndex: &focusedIndex,
		MainModel:    mainModel,
	}

	column := &(*model.calendar)[0]
	column.SetStyles(getActiveStyle())
	column.Focus()

	return model
}

func (model CalendarModel) prettyDateTime() string {
	return todayStyle.Render(
		fmt.Sprintf("Today is %s, %s %d",
			model.currentTime.Weekday().String(),
			model.currentTime.Month().String(),
			//[:3],
			model.currentTime.Day(),
		),
	)
}

func (model CalendarModel) Init() tea.Cmd {
	return nil
}

// Returns the row number of the current column
func (model CalendarModel) disableCurrent() int {
	column := &(*model.calendar)[*model.focusedIndex]
	column.SetStyles(getDefaultStyle())
	column.Blur()
	return column.Cursor()
}

// rowNum should be the row number of the previous
// column
func (model CalendarModel) enableCurrent(rowNum int) {
	column := &(*model.calendar)[*model.focusedIndex]
	column.SetStyles(getActiveStyle())
	column.Focus()
	column.SetCursor(rowNum)
}

func (model CalendarModel) MoveLeft() {
	rowNum := model.disableCurrent()
	*model.focusedIndex--
	model.enableCurrent(rowNum)
}

func (model CalendarModel) MoveRight() {
	rowNum := model.disableCurrent()
	*model.focusedIndex++
	model.enableCurrent(rowNum)
}

func (model CalendarModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	currentTime := time.Now()
	model.currentTime = &currentTime

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		modelStyle = modelStyle.
			Width(message.Width).
			Height(message.Height)

	case tea.KeyMsg:
		switch message.String() {

		case "alt+r":
			*model.targetMonth = model.currentTime.Month()
			*model.targetYear = model.currentTime.Year()

		case "left":
			if *model.focusedIndex == 0 {
				*model.targetMonth--
				if *model.targetMonth == 0 {
					*model.targetMonth = 12
					*model.targetYear--
				}

				rowNum := model.disableCurrent()
				*model.focusedIndex = 6
				model.enableCurrent(rowNum)

			} else if *model.focusedIndex > 0 {
				model.MoveLeft()
			}

		case "right":
			if *model.focusedIndex == 6 {
				*model.targetMonth++
				if *model.targetMonth == 13 {
					*model.targetMonth = 1
					*model.targetYear++
				}

				rowNum := model.disableCurrent()
				*model.focusedIndex = 0
				model.enableCurrent(rowNum)

			} else if *model.focusedIndex >= 0 && *model.focusedIndex < 6 {
				model.MoveRight()
			}

		case "tab":
			if *model.focusedIndex == -1 {
				*model.focusedIndex = 0
				(*model.calendar)[*model.focusedIndex].SetStyles(getActiveStyle())
				(*model.calendar)[*model.focusedIndex].Focus()
			} else {
				(*model.calendar)[*model.focusedIndex].SetStyles(getDefaultStyle())
				(*model.calendar)[*model.focusedIndex].Blur()
				*model.focusedIndex = -1
			}
		}
	}

	var command tea.Cmd
	if *model.focusedIndex == -1 {
		*model.events, command = model.events.Update(message)
	} else {
		(*model.calendar)[*model.focusedIndex],
			command = (*model.calendar)[*model.focusedIndex].Update(message)
	}

	return model, command
}

func (model CalendarModel) prettyTitle() string {
	return titleStyle.Render(
		fmt.Sprintf("%s %d",
			model.targetMonth.String(),
			*model.targetYear,
		),
	)
}

func (model CalendarModel) thisMonth() bool {
	return model.currentTime.Month() == *model.targetMonth &&
		model.currentTime.Year() == *model.targetYear
}

func (model CalendarModel) isToday(day int) bool {
	return model.thisMonth() && day == model.currentTime.Day()
}

func (model CalendarModel) calendarView() string {
	numDays := model.getNumDays()
	startOffset := model.getStartOffset()
	height := (numDays + startOffset + 6) / 7
	calendarArray := make([][]string, height)
	for i := range height {
		calendarArray[i] = make([]string, 7)
	}

	for i := 0; i < startOffset; i++ {
		calendarArray[0][i] = ""
	}
	x := 1
	for i := startOffset; i < 7; i++ {
		calendarArray[0][i] = fmt.Sprintf("%d", x)
		x++
	}
	for i := 1; i < height; i++ {
		for j := range 7 {
			if x <= numDays {
				calendarArray[i][j] = fmt.Sprintf("%d", x)
				x++
			} else {
				calendarArray[i][j] = ""
			}
		}
	}

	calendarViews := []string{}

	for i := range 7 {
		rows := []table.Row{}
		for j := range height {
			rows = append(rows,
				table.Row{calendarArray[j][i]},
			)
		}
		(*model.calendar)[i].SetRows(rows)
		calendarViews = append(calendarViews,
			(*model.calendar)[i].View(),
		)
	}
	text := strings.Builder{}
	text.WriteString(lipgloss.JoinHorizontal(lipgloss.Center,
		calendarViews...,
	) + "\n")

	return text.String()
}

func (model CalendarModel) View() string {
	text := strings.Builder{}
	text.WriteString(model.prettyDateTime() + "\n\n")
	text.WriteString(model.prettyTitle() + "\n")
	text.WriteString(model.calendarView() + "\n")
	eventStyle := lipgloss.NewStyle()
	text.WriteString(
		eventStyle.Render(model.events.View()) + "\n")
	return modelStyle.
		Background(lipgloss.Color("#ff0000")).
		Render(text.String())
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
