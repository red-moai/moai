package main

import (
	"github.com/Genekkion/moai/apps/bork"
	"github.com/Genekkion/moai/apps/calculator"
	"github.com/Genekkion/moai/apps/calendar"
	"github.com/Genekkion/moai/apps/diary"
	"github.com/Genekkion/moai/apps/todo"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModelInit func(external.MoaiModel) tea.Model

var (
	MOAI_APPS = map[string]ModelInit{
		"Bork":       bork.InitBork,
		"Calculator": calculator.InitCalculator,
		"Calendar":   calendar.InitCalendar,
		"Diary":      diary.InitDiary,
		"Todo":       todo.InitTodo,
	}

	AVAILABLE_APPS = []list.Item{
		MenuEntry{
			title:       "Bork",
			description: "A HTTP client for quick testing",
		},
		MenuEntry{
			title:       "Calendar",
			description: "Track your life",
		},
		MenuEntry{
			title:       "Calculator",
			description: "A simple calculator",
		},
		MenuEntry{
			title:       "Diary",
			description: "Your personal diary",
		},
		MenuEntry{
			title:       "Todo",
			description: "A simple todo list",
		},
	}
)

type MenuModel struct {
	list      list.Model
	mainModel external.MoaiModel
	style     lipgloss.Style
}

type MenuEntry struct {
	title       string
	description string
}

func (tabEntry MenuEntry) Title() string {
	return tabEntry.title
}
func (tabEntry MenuEntry) Description() string {
	return tabEntry.description
}
func (tabEntry MenuEntry) FilterValue() string {
	return tabEntry.title + " " + tabEntry.description
}

func InitMenu(mainModel external.MoaiModel) MenuModel {
	model := MenuModel{
		list: list.New(
			AVAILABLE_APPS,
			list.NewDefaultDelegate(),
			30,
			30,
		),
		mainModel: mainModel,
		style: lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()),
	}
	model.list.Title = "Available apps"

	return model
}

func (newTabModel MenuModel) Init() tea.Cmd {
	return nil
}

func (model *MenuModel) updateStyleDimensions(height int, width int) {
	model.style = model.style.
		Height(height).
		Width(width)
}

func (model MenuModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.list.SetWidth(message.Width)
		model.list.SetHeight(message.Height)
		model.updateStyleDimensions(message.Height, message.Width)

		return model, nil
	case tea.KeyMsg:

	}

	var command tea.Cmd
	model.list, command = model.list.Update(message)
	return model, command
}

func (model MenuModel) View() string {
	return model.style.Render(model.list.View())
}
