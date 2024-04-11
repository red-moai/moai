package components

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	defaultListStyle = lipgloss.NewStyle().
		Margin(1, 2)
)

type DefaultListModel struct {
	List         list.Model
	SelectedItem list.Item
	CustomStyle  *lipgloss.Style
}

func InitDefaultList(items []list.Item, title string,
	width int, height int, customStyle *lipgloss.Style,
	keyBindings ...key.Binding) DefaultListModel {
	model := DefaultListModel{
		List: list.New(
			items,
			list.NewDefaultDelegate(),
			width,
			height,
		),
		SelectedItem: nil,
		CustomStyle:  customStyle,
	}
	model.List.Title = title
	model.List.Help.ShowAll = false
	model.List.KeyMap.Quit.Unbind()
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
	if model.CustomStyle != nil {
		return model.CustomStyle.Render(model.List.View())
	}
	return defaultListStyle.Render(model.List.View())
}