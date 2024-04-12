package home

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Genekkion/moai/external"
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
			Align(lipgloss.Center)
)

type HomeModel struct {
	HomeChoices  []string
	HomeCursor   int
	HomeSelected map[int]struct{}
	HomeQuote    string
	MainModel    external.MoaiModel
}

func InitHome(mainModel external.MoaiModel) tea.Model {
	model := HomeModel{
		HomeChoices: []string{

			"Oonga boonga",
			"boo ya",
		},
		HomeSelected: map[int]struct{}{},
		HomeQuote:    homeQuotes[rand.Intn(len(homeQuotes))],
		MainModel:    mainModel,
	}

	return model
}

func (model HomeModel) Init() tea.Cmd {
	return nil
}

func (model HomeModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {

	case tea.KeyMsg:
		switch message.String() {

		case "up", "k":
			if model.HomeCursor > 0 {
				model.HomeCursor--
			}

		case "down", "j":
			if model.HomeCursor < len(model.HomeChoices)-1 {
				model.HomeCursor++
			}

		case "enter", " ":
			_, ok := model.HomeSelected[model.HomeCursor]
			if ok {
				delete(model.HomeSelected, model.HomeCursor)
			} else {
				model.HomeSelected[model.HomeCursor] = struct{}{}
			}
		}
	}

	return model, nil
}

func (model HomeModel) View() string {

	if !model.MainModel.IsReady() {
		return modelStyle.
			Width(model.MainModel.AvailableWidth()).
			Height(model.MainModel.AvailableHeight()).
			Render("booya")
	}

	text := strings.Builder{}

	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(0, 1).
		Width(
			min(
				50,
				lipgloss.Width(model.HomeQuote)+2,
			),
		)

	// Header
	text.WriteString(headerStyle.Render(model.HomeQuote) + "\n")

	// Iterate over our choices
	for i, choice := range model.HomeChoices {

		cursor := " " // no cursor
		if model.HomeCursor == i {
			cursor = ">" // cursor!
		}

		checked := " " // not selected
		if _, ok := model.HomeSelected[i]; ok {
			checked = "x" // selected!
		}

		text.WriteString(
			fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice),
		)
	}
	return modelStyle.
		Width(model.MainModel.AvailableWidth()).
		Height(model.MainModel.AvailableHeight()).
		Render(text.String())
}
