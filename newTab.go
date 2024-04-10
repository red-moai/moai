package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelInit func() tea.Model

var (
	MOAI_MODELS = map[string]ModelInit{
		"Home": func() tea.Model {
			return InitHome()
		},
		"Diary": func() tea.Model {
			return InitDiary()
		},
	}

	AVAILABLE_MODELS = []list.Item{
	}
)

type NewTabModel struct {
	SearchInput   textinput.Model
	SearchDisplay string

	ModelList list.Model
}

func InitNewTab() NewTabModel {
	model := NewTabModel{
		SearchInput:   textinput.New(),
		SearchDisplay: "",
	}

	model.SearchInput.Placeholder = "Search"
	model.SearchInput.Focus()
	model.SearchInput.CharLimit = 255
	model.SearchInput.Width = 20

	return model
}

func (newTabModel NewTabModel) Init() tea.Cmd {
	return nil
}

func (model NewTabModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:

		switch message.Type {
		case tea.KeyEnter:
			model.SearchDisplay = model.SearchInput.Value()
		}

	case error:
		model.SearchDisplay = message.Error()
		return model, nil
	}

	var command tea.Cmd
	model.SearchInput, command = model.SearchInput.Update(message)
	return model, command
}

func (model NewTabModel) View() string {

	text := model.SearchInput.View() + "\n"
	text += model.SearchDisplay + "\n"

	return text
}
