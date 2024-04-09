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

func (model *Model) initDiary() {
	textInput := textinput.New()
	textInput.Placeholder = "Search"
	textInput.Focus()
	textInput.CharLimit = 255
	textInput.Width = 20

	model.diarySearch = textInput
}

func (model *Model) updateDiary(message tea.Msg) (tea.Model, tea.Cmd) {
	var command tea.Cmd
	switch message := message.(type) {
	case tea.KeyMsg:

		switch message.Type {
		case tea.KeyEnter:
			model.diaryDisplay = model.diarySearch.Value()
		}

	case error:
		model.Error = message
		return model, nil
	}

	model.diarySearch, command = model.diarySearch.Update(message)
	return model, command
}

func (model *Model) renderDiary() string {
	text := model.diarySearch.View() + "\n"
	text += model.diaryDisplay + "\n"
	return text
}
