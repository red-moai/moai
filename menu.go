package main

import (
	"github.com/Genekkion/moai/apps/bork"
	"github.com/Genekkion/moai/apps/calculator"
	"github.com/Genekkion/moai/apps/calendar"
	"github.com/Genekkion/moai/apps/diary"
	"github.com/Genekkion/moai/apps/gpt"
	"github.com/Genekkion/moai/apps/todo"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	MOAI_APPS = []list.Item{
		MenuEntry{
			title:       "Bork",
			description: "A HTTP client for quick testing",
			initialiser: bork.InitBork,
		},
		MenuEntry{
			title:       "Calendar",
			description: "Track your life",
			initialiser: calendar.InitCalendar,
		},
		MenuEntry{
			title:       "Calculator",
			description: "A simple calculator",
			initialiser: calculator.InitCalculator,
		},
		MenuEntry{
			title:       "Diary",
			description: "Your personal diary",
			initialiser: diary.InitDiary,
		},
		MenuEntry{
			title:       "GPT",
			description: "Access OpenAI's models",
			initialiser: gpt.InitGPT,
		},
		MenuEntry{
			title:       "Todo",
			description: "A simple todo list",
			initialiser: todo.InitTodo,
		},
	}
)

type MenuModel struct {
	list         list.Model
	table        table.Model
	tableColumns []table.Column
	keymap       MenuKeyMap
	listFocused  bool
	showHelp     bool
}

type ModelInit func(external.MoaiModel) tea.Model

type MenuEntry struct {
	title       string
	description string
	initialiser ModelInit
}

func (entry MenuEntry) Title() string {
	return entry.title
}
func (entry MenuEntry) Description() string {
	return entry.description
}
func (entry MenuEntry) FilterValue() string {
	return entry.title + " " + entry.description
}

func extractTabs(entries TabEntries) []table.Row {
	numEntries := entries.Len() - 1
	if numEntries == 0 {
		return []table.Row{}
	}
	rows := make([]table.Row, numEntries)
	validEntries := entries[1:]
	for i := range numEntries {
		rows[i] = []string{
			validEntries[i].title,
			validEntries[i].timeElapsedString(),
		}
	}
	return rows
}

func InitMenu(mainModel Model) tea.Model {
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
			MOAI_APPS,
			list.NewDefaultDelegate(),
			50,
			50,
		),
		table: table.New(
			table.WithColumns(recentlyUsedColumns),
			table.WithRows(extractTabs(mainModel.tabs)),
			table.WithFocused(false),
			table.WithStyles(tableStyle),
		),
		tableColumns: recentlyUsedColumns,

		keymap:      initMenuKeyMap(mainModel.ModKey()),
		showHelp:    true,
		listFocused: true,
	}

	model.list.Title = "Available apps"
	model.list.SetShowHelp(false)
	model.list.KeyMap.ForceQuit.Unbind()
	model.list.KeyMap.Quit.Unbind()
	model.list.KeyMap.ShowFullHelp.Unbind()
	model.list.KeyMap.CloseFullHelp.Unbind()

	model.table.GotoTop()
	model.updateDimensions(mainModel.latestWindowMsg)

	return model
}

func (newTabModel MenuModel) Init() tea.Cmd {
	return nil
}

var (
	menuModelStyle = lipgloss.NewStyle().
			Padding(1).
			AlignHorizontal(lipgloss.Center).
			Border(lipgloss.RoundedBorder())
	menuWidgetStyle = lipgloss.NewStyle()
	menuHelpStyle   = lipgloss.NewStyle()
)

func (model *MenuModel) updateDimensions(message tea.Msg) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		menuModelStyle = menuModelStyle.
			Height(message.Height - 2).
			Width(message.Width - 2)
		menuHelpStyle = menuHelpStyle.
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

		menuWidgetStyle = menuWidgetStyle.
			Height(widgetHeight - 14).
			Width(widgetWidth + 1)
	}
}

func (model MenuModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.updateDimensions(message)

		return model, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(message, model.keymap.Help):
			model.showHelp = !model.showHelp
			return model, nil
		case key.Matches(message, model.keymap.Focus):
			if model.listFocused {
				model.table.Focus()
			} else {
				model.table.Blur()
			}
			model.listFocused = !model.listFocused
			return model, nil
		}
	}

	var command tea.Cmd
	if model.listFocused {
		switch message := message.(type) {

		case tea.KeyMsg:
			switch message.String() {
			case "enter":
				return model, func() tea.Msg {
					return model.list.SelectedItem().(MenuEntry)
				}
			}

			switch {
			case key.Matches(message, model.keymap.Home):
				return model, func() tea.Msg {
					return "switchHome"
				}
			}
		}

		model.list, command = model.list.Update(message)
	} else {
		switch message := message.(type) {

		case tea.KeyMsg:
			switch message.String() {
			case "enter":
				return model, func() tea.Msg {
					return SetIndexMessage{
						index: model.table.Cursor() + 1,
					}
				}
			}
		}
		model.table, command = model.table.Update(message)
	}

	return model, command
}

type SetIndexMessage struct {
	index int
}

func (model MenuModel) listView() string {
	return menuWidgetStyle.Render(model.list.View())
}

func (model MenuModel) tableView() string {
	return menuWidgetStyle.Render(model.table.View())
}

func (model MenuModel) helpView() string {
	return menuHelpStyle.Render(
		model.list.Help.View(model.keymap),
	)
}

func (model MenuModel) View() string {
	gap := ""
	if menuModelStyle.GetWidth()%2 != 0 {
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
	return menuModelStyle.Render(text)
}
