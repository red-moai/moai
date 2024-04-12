package home

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	homeQuotes = []string{
		"Yesterday is history, tomorrow is a mystery, but today is a gift. That is why it is called the present.",
		"The fellas",
	}

	modelStyle = lipgloss.NewStyle().
			Padding(1).
		// Background(lipgloss.Color("#24283B")).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder())

	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(0, 1).
			Margin(1, 0)

	USERNAME      = "Mr Moai"
	PRETTY_LAYOUT = "02 January 2006, 15:04:05"
)

type HomeModel struct {
	quote       string
	spinner     spinner.Model
	currentTime time.Time
	helpModel   help.Model
	keyMap      KeyMap
	showHelp    bool

	MainModel external.MoaiModel
}

func (model HomeModel) GetSpinner() spinner.Model {
	return model.spinner
}

func InitHome(mainModel external.MoaiModel) tea.Model {
	model := HomeModel{
		quote: homeQuotes[rand.Intn(len(homeQuotes))],
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Moon),
		),
		currentTime: time.Now(),
		helpModel:   help.New(),
		keyMap:      initKeyMap(mainModel.ModKey()),
		MainModel:   mainModel,
	}

	return model
}

func (model HomeModel) Init() tea.Cmd {
	return model.spinner.Tick
}

func (model *HomeModel) updateTime() {
	model.currentTime = time.Now()
}

func (model HomeModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		modelStyle = modelStyle.
			Width(message.Width - 4).
			Height(message.Height - 2)

		model.helpModel.Width = message.Width

	case spinner.TickMsg:
		model.updateTime()

		var command tea.Cmd
		model.spinner, command = model.spinner.Update(message)
		return model, command

	case tea.KeyMsg:
		switch {
		case key.Matches(message, model.keyMap.Help):
			model.toggleHelp()
		}

	}

	return model, nil
}

func (model *HomeModel) toggleHelp() {
	model.showHelp = !model.showHelp
}

func (model HomeModel) welcomeView() string {
	welcomeStyle := lipgloss.NewStyle().
		Margin(1)
	text := strings.Builder{}

	hour := model.currentTime.Hour()
	if hour >= 0 && hour < 6 {
		text.WriteString("Good night")
	} else if hour < 12 {
		text.WriteString("Good morning")
	} else if hour < 18 {
		text.WriteString("Good afternoon")
	} else {
		text.WriteString("Good evening")
	}

	text.WriteString(", ")
	text.WriteString(USERNAME)
	text.WriteString(" ")
	text.WriteString(model.spinner.View())
	return welcomeStyle.Render(text.String())
}

func (model HomeModel) quoteView() string {
	text := "Quote of the day:\n"

	return headerStyle.
		Width(
			min(
				max(
					lipgloss.Width(model.quote)+4,
					lipgloss.Width(text)+4,
				),
				modelStyle.GetWidth()/2,
			),
		).
		Render(text + model.quote)
}

func (model HomeModel) timeView() string {
	date := model.currentTime.Format(PRETTY_LAYOUT)
	day, _ := strconv.Atoi(date[:2])
	suffix := "th"
	switch day {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	text := fmt.Sprintf("%d%s%s", day, suffix, date[2:])

	return "It is currently: " + text
}

func (model HomeModel) View() string {

	text := strings.Builder{}

	// Header
	text.WriteString(model.quoteView() + "\n")

	text.WriteString(model.timeView() + "\n")

	text.WriteString(model.welcomeView() + "\n")

	if model.showHelp {
		text.WriteString(model.helpModel.View(model.keyMap))
	}

	return modelStyle.
		Render(text.String())
}
