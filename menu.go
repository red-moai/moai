package main

import (
	"github.com/Genekkion/moai/apps/bork"
	"github.com/Genekkion/moai/apps/calculator"
	"github.com/Genekkion/moai/apps/calendar"
	"github.com/Genekkion/moai/apps/diary"
	"github.com/Genekkion/moai/apps/todo"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
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

	fakeRecentlyUsed = []table.Row{
		{"Bork", "2 mins ago"},
		{"Notes", "1 hour ago"},
	}
)

type MenuModel struct {
	list         list.Model
	table        table.Model
	tableColumns []table.Column
	listStyle    lipgloss.Style
	modelStyle   lipgloss.Style
	helpStyle    lipgloss.Style
	keymap       MenuKeyMap
	listFocused  bool
	showHelp     bool

	mainModel external.MoaiModel
}

type MenuEntry struct {
	title       string
	description string
}

func (menuEntry MenuEntry) Title() string {
	return menuEntry.title
}
func (menuEntry MenuEntry) Description() string {
	return menuEntry.description
}
func (menuEntry MenuEntry) FilterValue() string {
	return menuEntry.title + " " + menuEntry.description
}

func InitMenu(mainModel external.MoaiModel) MenuModel {
	recentlyUsedColumns := []table.Column{
		{Title: "Application", Width: 15},
		{Title: "Last used", Width: 10},
	}

	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		Foreground(lipgloss.Color("#B4F9F8")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true)

	model := MenuModel{
		list: list.New(
			AVAILABLE_APPS,
			list.NewDefaultDelegate(),
			50,
			50,
		),
		table: table.New(
			table.WithColumns(recentlyUsedColumns),
			table.WithRows(fakeRecentlyUsed),
			table.WithFocused(false),
			table.WithStyles(tableStyle),
		),
		tableColumns: recentlyUsedColumns,

		mainModel: mainModel,

		modelStyle: lipgloss.NewStyle().
			Padding(1).
			AlignHorizontal(lipgloss.Center).
			Border(lipgloss.RoundedBorder()),
		listStyle: lipgloss.NewStyle(),
			//Border(lipgloss.RoundedBorder()),
		helpStyle: lipgloss.NewStyle(),

		keymap:      initMenuKeyMap(mainModel.ModKey()),
		showHelp:    true,
		listFocused: true,
	}
	model.list.Title = "Available apps"
	model.list.SetShowHelp(false)
	model.list.DisableQuitKeybindings()
	model.list.KeyMap.ShowFullHelp.Unbind()
	model.list.KeyMap.CloseFullHelp.Unbind()

	model.table.GotoTop()

	return model
}

func (newTabModel MenuModel) Init() tea.Cmd {
	return nil
}

func (model MenuModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.modelStyle = model.modelStyle.
			//Background(lipgloss.Color("#00ff00")).
			Height(message.Height - 2).
			Width(message.Width - 2)
		model.helpStyle = model.helpStyle.
			Width(message.Width - 4)

		widgetHeight := message.Height - 6 -
			lipgloss.Height(model.helpView())
		widgetWidth := (message.Width-4)/2 - 4

		model.list.SetHeight(widgetHeight)
		model.list.SetWidth(widgetWidth)
		model.table.SetHeight(widgetHeight - 2)

		newColumns := make([]table.Column, len(model.tableColumns))

		for i := range len(model.tableColumns) {
			newColumns[i].Title = model.tableColumns[i].Title
			newColumns[i].Width = ((message.Width-4)/2 - 7) / 2
		}
		model.table.SetColumns(newColumns)

		model.listStyle = model.listStyle.
			Height(widgetHeight - 14).
			Width(widgetWidth + 1)

		return model, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(message, model.keymap.Exit):
			model.mainModel.ToggleMenu()

		case key.Matches(message, model.keymap.Help):
			model.showHelp = !model.showHelp

		case key.Matches(message, model.keymap.Focus):
			if model.listFocused {
				model.table.Focus()
			} else {
				model.table.Blur()
			}
			model.listFocused = !model.listFocused
		}
	}

	var command tea.Cmd
	if model.listFocused {
		model.list, command = model.list.Update(message)
	} else {
		model.table, command = model.table.Update(message)
	}

	return model, command
}

func (model MenuModel) listView() string {
	return model.listStyle.Render(model.list.View())
}

func (model MenuModel) tableView() string {
	return model.listStyle.Render(model.table.View())
}

func (model MenuModel) helpView() string {
	return model.helpStyle.Render(
		model.list.Help.View(model.keymap),
	)
}

func (model MenuModel) View() string {
	gap := ""
	if model.modelStyle.GetWidth()%2 != 0 {
		gap = " "
	}
	text := lipgloss.JoinHorizontal(
		lipgloss.Center,
		model.listView(),
		gap,
		model.tableView(),
	)
	if model.showHelp {
		text = lipgloss.JoinVertical(
			lipgloss.Center,
			text,
			model.helpView(),
		)
	}
	return model.modelStyle.Render(text)
}
