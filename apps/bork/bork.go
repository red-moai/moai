package bork

import (
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BorkModel struct {
	borkCategories []string
	keymap         KeyMap

	list list.Model
}

var (
	borkCategories = []string{
		"ChatsTodo",
		"PottySense",
		"New category",
	}

	modelStyle = lipgloss.NewStyle().
			Padding(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#CFC9C2"))
)

func InitBork(mainModel external.MoaiModel) tea.Model {
	model := BorkModel{
		list: list.New(
			borkEntries,
			list.NewDefaultDelegate(),
			30,
			30,
		),
		keymap:         initKeyMap(mainModel.ModKey()),
		borkCategories: borkCategories,
	}
	model.list.Title = "Bork Bork! 🐶"

	model.updateDimensions(mainModel.GetLatestWindowMessage())
	model.list.SetShowHelp(false)

	return model
}

func (model BorkModel) Init() tea.Cmd {
	return nil
}

func (model *BorkModel) updateDimensions(message tea.Msg) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		newHeight := message.Height - 2
		newWidth := message.Width - 2
		modelStyle = modelStyle.
			Height(newHeight).
			Width(newWidth)
		model.list.SetHeight(newHeight - 2)
		model.list.SetWidth(newWidth - 2)

	}
}

func (model BorkModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.updateDimensions(message)
		return model, nil

	case tea.KeyMsg:

		if key.Matches(message, model.keymap.AddNew) {
			borkFormModel := initBorkForm(&model)

			return borkFormModel, nil
		}
	}

	var command tea.Cmd
	model.list, command = model.list.Update(message)
	return model, command
}

func (model BorkModel) View() string {
	return modelStyle.Render(model.list.View())
}
