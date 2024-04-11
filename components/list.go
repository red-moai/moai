package components

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	fullListStyle = lipgloss.NewStyle().
			Padding(1, 2)
	fullListTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#25A065")).
				Padding(0, 1)
	fullListStatusStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type FullListModel struct {
	List         list.Model
	SelectedItem list.Item
	CustomStyle  *lipgloss.Style
}

func InitFullList(items []list.Item, title string,
	width int, height int, customStyle *lipgloss.Style,
	keyBindings ...key.Binding) FullListModel {
	

	model := FullListModel{
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
	return model
}

func (model FullListModel) Update(message tea.Msg) (FullListModel, tea.Cmd) {

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

func (model FullListModel) View() string {
	if model.CustomStyle != nil {
		return model.CustomStyle.Render(model.List.View())
	}
	return defaultListStyle.Render(model.List.View())
}
