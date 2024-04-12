package components

import (
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	defaultListStyle = lipgloss.NewStyle()
)

type DefaultListModel struct {
	List         list.Model
	SelectedItem list.Item
	CustomStyle  *lipgloss.Style
	MainModel    external.MoaiModel
}

func InitDefaultList(items []list.Item, title string,
	mainModel external.MoaiModel, listStyle *list.Styles,
	delegateStyles *list.DefaultItemStyles,
	keyBindings ...key.Binding) DefaultListModel {

	model := DefaultListModel{
		List: list.New(
			items,
			list.NewDefaultDelegate(),
			30,
			30,
		),
		SelectedItem: nil,
		MainModel:    mainModel,
	}
	model.List.Title = title
	model.List.Help.ShowAll = false
	model.List.KeyMap.Quit.Unbind()
	if listStyle != nil {
		model.List.Styles = *listStyle
	}

	if delegateStyles != nil {
		delegate := list.NewDefaultDelegate()
		delegate.Styles = *delegateStyles
		model.List.SetDelegate(delegate)
	}

	/*
		model.List.AdditionalShortHelpKeys = func() []key.Binding {
			return keyBindings
		}
	*/
	model.List.AdditionalFullHelpKeys = func() []key.Binding {
		return keyBindings
	}
	return model
}

func (model DefaultListModel) Update(message tea.Msg) (DefaultListModel, tea.Cmd) {

	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.Type {
		case tea.KeyEnter:
			model.SelectedItem = model.List.SelectedItem()
		}

	case tea.WindowSizeMsg:
		x, y := defaultListStyle.GetFrameSize()
		model.List.SetSize(message.Width-x, message.Height-y)
	}

	var command tea.Cmd
	model.List, command = model.List.Update(message)
	return model, command
}

func (model DefaultListModel) View() string {
	return lipgloss.NewStyle().
		Render(model.List.View())

	//return model.List.View()
}
