package main

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type Diary struct {
	Id        int
	Title     string
	Text      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type DiaryModel struct {
	SearchInput   textinput.Model
	SearchDisplay string
}

func mockDiaryData() []Diary {
	timeNow := time.Now()
	return []Diary{
		{
			Id:        0,
			Title:     "Entry 1",
			Text:      "diary text 1",
			CreatedAt: &timeNow,
			UpdatedAt: &timeNow,
		},
		{
			Id:        1,
			Title:     "Entry 2",
			Text:      "diary text 2",
			CreatedAt: &timeNow,
			UpdatedAt: &timeNow,
		},
	}
}

func InitDiary() DiaryModel {
	model := DiaryModel{
		SearchInput:   textinput.New(),
		SearchDisplay: "",
	}
	model.SearchInput.Placeholder = "Search"
	model.SearchInput.Focus()
	model.SearchInput.CharLimit = 255
	model.SearchInput.Width = 20

	return model
}

func (model DiaryModel) Init() tea.Cmd {
	return nil
}

func (model DiaryModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var command tea.Cmd
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

	model.SearchInput, command = model.SearchInput.Update(message)
	return model, command
}

func (model DiaryModel) View() string {
	text := "Diary\n"
	text += model.SearchInput.View() + "\n"
	text += model.SearchDisplay + "\n"
	return text
}
