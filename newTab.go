package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type NewTabModel struct{}

func (model Model) initNewTab() {
	textInput := textinput.New()
	textInput.Placeholder = "Search"
	textInput.Focus()
	textInput.CharLimit = 255
	textInput.Width = 20

	model.newTabSearch = textInput
}

func (newTabModel NewTabModel) Init() tea.Cmd {
	return nil
}

func (newTabModel NewTabModel) Update(model *Model, message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:

		switch message.Type {
		case tea.KeyEnter:
			model.newTabDisplay = model.newTabSearch.Value()
		}

	case error:
		model.Error = message
		return model, nil
	}

	var command tea.Cmd
	model.newTabSearch, command = model.newTabSearch.Update(message)
	return model, command
}

func (newTabModel NewTabModel) View(model Model) string {
	text := "new Tab\n"
	text += model.newTabSearch.View() + "\n"
	text += model.newTabDisplay + "\n"
	return text
}
