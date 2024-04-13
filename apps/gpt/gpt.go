package gpt

import (
	"strings"

	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	modelStyle = lipgloss.NewStyle().
			Padding(1).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder())

	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(0, 1).
			Margin(1, 0)
)

type GPTModel struct {
	spinner   spinner.Model
	helpModel help.Model
	keyMap    KeyMap
	showHelp  bool

	MainModel external.MoaiModel
}

func (model GPTModel) GetSpinner() spinner.Model {
	return model.spinner
}

func InitGPT(mainModel external.MoaiModel) tea.Model {
	model := GPTModel{
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Moon),
		),
		helpModel: help.New(),
		keyMap:    initKeyMap(mainModel.ModKey()),
		MainModel: mainModel,
	}

	return model
}

func (model GPTModel) Init() tea.Cmd {
	return model.spinner.Tick
}

func (model GPTModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		modelStyle = modelStyle.
			Width(message.Width - 2).
			Height(message.Height - 2)

		model.helpModel.Width = message.Width

	case spinner.TickMsg:

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

func (model *GPTModel) toggleHelp() {
	model.showHelp = !model.showHelp
}

func (model GPTModel) View() string {

	text := strings.Builder{}

	text.WriteString("hi")

	return modelStyle.
		Render(text.String())
}
