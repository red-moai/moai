package home

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Genekkion/moai/external"
	"github.com/Genekkion/moai/log"
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
			Align(lipgloss.Center)

	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(0, 1).
			MarginTop(1)

	quoteStyle = lipgloss.NewStyle()

	welcomeStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			Margin(1, 0).
			Bold(true)

	PRETTY_LAYOUT = "02 January 2006, 3:04:05pm"
)

type HomeModel struct {
	quote         string
	spinner       spinner.Model
	currentTime   time.Time
	keyMap        KeyMap
	username      string
	quoteReceived bool
	debug         string

	MainModel *external.MoaiModel
}

func (model HomeModel) GetSpinner() spinner.Model {
	return model.spinner
}

var MoaiSpinner = spinner.Spinner{
	Frames: []string{
		"🗿", "🌴", "🌊", "🏝️", "🌺", "🌋", "🌄", "🚣", "🏄", "🎣", "🍹", "🌅", "🌞", "🌇", "🌉", "🌌", "🌠", "🌟", "🌙", "🌃", "🎆",
	},
	FPS: time.Second / 3,
}

type QuoteMessage struct {
	quote string
	err   error
}

func getQuote() tea.Msg {
	response, err := http.DefaultClient.Get("http://localhost:3000/quote")
	if err != nil || response.StatusCode != http.StatusOK {
		log.ErrorWrapper(err)
		return QuoteMessage{
			quote: homeQuotes[rand.Intn(len(homeQuotes))],
			err:   nil,
		}
	}
	defer response.Body.Close()

	type ResponseBody struct {
		Id    int    `json:"id"`
		Quote string `json:"quote"`
	}

	var responseBody ResponseBody
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		log.ErrorWrapper(err)
		return QuoteMessage{
			quote: homeQuotes[rand.Intn(len(homeQuotes))],
			err:   nil,
		}
	}

	return QuoteMessage{
		quote: responseBody.Quote,
		err:   nil,
	}
}

func InitHome(mainModel external.MoaiModel) tea.Model {
	model := HomeModel{
		spinner: spinner.New(
			//spinner.WithSpinner(MoaiSpinner),
			spinner.WithSpinner(spinner.Dot),
		),
		currentTime:   time.Now(),
		keyMap:        initKeyMap(mainModel.ModKey()),
		MainModel:     &mainModel,
		username:      mainModel.Username(),
		quoteReceived: false,
		quote:         "Loading...",
	}
	colorScheme := mainModel.ColorScheme()
	modelStyle = modelStyle.
		Background(colorScheme.BG3())
	headerStyle = headerStyle.
		Foreground(colorScheme.FG3()).
		BorderBackground(colorScheme.BG3()).
		Background(colorScheme.BG3()).
		BorderForeground(colorScheme.FG1())
	quoteStyle = quoteStyle.
		Foreground(colorScheme.FG2())
	welcomeStyle = welcomeStyle.
		Background(colorScheme.BG3())

	return model
}

func (model HomeModel) Init() tea.Cmd {
	return tea.Batch(model.spinner.Tick, getQuote)
}

func (model *HomeModel) updateTime() {
	model.currentTime = time.Now()
}

func (model HomeModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		modelStyle = modelStyle.
			Width(message.Width).
			Height(message.Height - 2)

	case QuoteMessage:
		model.quoteReceived = true
		if message.err != nil {
			return model, nil
		}
		model.quote = message.quote
		return model, nil

	case spinner.TickMsg:
		model.updateTime()

		var command tea.Cmd
		model.spinner, command = model.spinner.Update(message)
		return model, command

	case tea.KeyMsg:
		model.debug = message.String()
	}

	return model, nil
}

func (model HomeModel) welcomeView() string {
	text := strings.Builder{}
	emojis := []string{"🌅", "🌞", "🌆", "🌙"}
	var index int

	hour := model.currentTime.Hour()
	if hour >= 0 && hour < 6 {
		text.WriteString("Good night")
		index = 3
	} else if hour < 12 {
		text.WriteString("Good morning")
		index = 0
	} else if hour < 18 {
		text.WriteString("Good afternoon")
		index = 1
	} else {
		text.WriteString("Good evening")
		index = 2
	}

	text.WriteString(", ")
	text.WriteString(model.username)
	text.WriteString(" ")
	//text.WriteString(model.spinner.View())
	text.WriteString(emojis[index])
	return welcomeStyle.Render(text.String())
}

func (model HomeModel) quoteView() string {
	var text string
	if !model.quoteReceived {
		text = fmt.Sprintf(
			"%s\n%s%s",
			quoteStyle.Render("Quote of the day:"),
			"loading ",
			model.spinner.View(),
		)

	} else {
		text = fmt.Sprintf(
			"%s\n%s",
			quoteStyle.Render("Quote of the day:"),
			model.quote,
		)
	}

	return headerStyle.
		Width(
			min(
				lipgloss.Width(text)+4,
				modelStyle.GetWidth()/2,
			),
		).
		Render(text)
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

	text.WriteString(model.quoteView())
	text.WriteString("\n")
	text.WriteString(model.welcomeView())
	text.WriteString("\n")
	text.WriteString(model.timeView())
	text.WriteString("\n")
	text.WriteString(fmt.Sprintf("Debug: %s", model.debug))

	return modelStyle.Render(text.String()) 
}
