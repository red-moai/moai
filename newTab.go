package main

import (
	"github.com/Genekkion/moai/apps/bork"
	"github.com/Genekkion/moai/apps/calculator"
	"github.com/Genekkion/moai/apps/calendar"
	"github.com/Genekkion/moai/apps/diary"
	"github.com/Genekkion/moai/apps/home"
	"github.com/Genekkion/moai/apps/todo"
	"github.com/Genekkion/moai/components"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModelInit func(external.MoaiModel) tea.Model

type SwapModelFunc func(string, tea.Model)

var (
	MOAI_APPS = map[string]ModelInit{
		"Home":       home.InitHome,
		"Bork":       bork.InitBork,
		"Calculator": calculator.InitCalculator,
		"Calendar":   calendar.InitCalendar,
		"Diary":      diary.InitDiary,
		"Todo":       todo.InitTodo,
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
		TabEntry{
			title:       "Todo",
			description: "A simple todo list",
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
			mainModel,
			func() *list.Styles {
				styles := list.DefaultStyles()
				styles.Title = styles.Title.
					Foreground(lipgloss.Color("#7AA2F7")).
					Background(lipgloss.NoColor{}).
					Bold(true)

				return &styles
			}(),
			func() *list.DefaultItemStyles {
				delegate := list.NewDefaultItemStyles()
				delegate.SelectedTitle = delegate.SelectedTitle.
					Foreground(lipgloss.Color("#F7768E"))

				return &delegate
			}(),
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
		return replacementModel, command

	}
	return model, command
}

func (model NewTabModel) View() string {
	return lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		//Background(lipgloss.Color("#00FF00")).
		Render(model.ModelList.View())
}
