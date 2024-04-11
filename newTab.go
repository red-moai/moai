package main

import (
	"github.com/Genekkion/moai/apps/bork"
	"github.com/Genekkion/moai/apps/calculator"
	"github.com/Genekkion/moai/apps/calendar"
	"github.com/Genekkion/moai/apps/diary"
	"github.com/Genekkion/moai/apps/home"
	"github.com/Genekkion/moai/components"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelInit func(external.MoaiModel) tea.Model

type SwapModelFunc func(string, tea.Model)

var (
	MOAI_APPS = map[string]ModelInit{
		"Home":       home.InitHome,
		"Diary":      diary.InitDiary,
		"Bork":       bork.InitBork,
		"Calculator": calculator.InitCalculator,
		"Calendar":   calendar.InitCalendar,
	}

	AVAILABLE_APPS = []list.Item{
		TabEntry{
			title:       "Home",
			description: "Dashboard",
		},
		TabEntry{
			title:       "Bork",
			description: "A HTTP client for quick testing",
		},
		TabEntry{
			title:       "Calendar",
			description: "Track your life",
		},
		TabEntry{
			title:       "Calculator",
			description: "A simple calculator",
		},
		TabEntry{
			title:       "Diary",
			description: "Your personal diary",
		},
	}
)

type NewTabModel struct {
	ModelList components.DefaultListModel
	mainModel external.MoaiModel
}

type TabEntry struct {
	title       string
	description string
}

func (tabEntry TabEntry) Title() string {
	return tabEntry.title
}
func (tabEntry TabEntry) Description() string {
	return tabEntry.description
}
func (tabEntry TabEntry) FilterValue() string {
	return tabEntry.title + " " + tabEntry.description
}

func InitNewTab(mainModel external.MoaiModel) NewTabModel {

	model := NewTabModel{
		ModelList: components.InitDefaultList(
			AVAILABLE_APPS,
			"Available apps",
			30,
			20,
			nil,
		),
		mainModel: mainModel,
	}
	return model
}

func (newTabModel NewTabModel) Init() tea.Cmd {
	return nil
}

func (model NewTabModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var command tea.Cmd
	model.ModelList, command = model.ModelList.Update(message)

	// If selected item, swap with current tab
	if model.ModelList.SelectedItem != nil {
		title := model.ModelList.SelectedItem.(TabEntry).title
		replacementModel := MOAI_APPS[title](model.mainModel)

		model.mainModel.SwapActiveModel(title, replacementModel)
		return replacementModel, command

	}
	return model, command
}

func (model NewTabModel) View() string {
	return model.ModelList.View() + "\n"
}
