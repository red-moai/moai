package main

import (
	"time"

	"github.com/Genekkion/moai/apps/home"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	// Globals
	Error  error
	modkey string

	ActiveTab   int
	PreviousTab int
	tabs        TabEntries

	keyMap      GlobalKeyMap
	onHome      bool
	menuSpawned bool

	latestWindowMsg tea.Msg
}

// Initialises the model to be ran by bubbletea
func InitModel() Model {
	model := Model{
		onHome:      true,
		menuSpawned: false,
		ActiveTab:   0,
		PreviousTab: 0,
		modkey:      getModkey(),
	}

	model.keyMap = initGlobalKeyMap(model.modkey)
	model.tabs = TabEntries{
		{
			title: "Home",
			model: home.InitHome(&model).(home.HomeModel),
		},
	}

	return model
}

func (model Model) ModKey() string {
	return model.modkey
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("M O A I 🗿"),
		model.tabs[0].model.Init(),
	)
}

func (model *Model) spawnMenu() {
	model.PreviousTab = model.ActiveTab
	model.ActiveTab = len(model.tabs)
	model.tabs = append(model.tabs,
		TabEntry{
			title: "New Tab",
			model: InitMenu(*model),
		},
	)
	model.menuSpawned = true
}

// Main function to update contents of application.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.tabs[0].model, _ = model.tabs[0].model.Update(message)

		model.latestWindowMsg = message

	case tea.KeyMsg:
		keypress := message.String()

		switch keypress {
		case "ctrl+c":
			return model, tea.Quit
		}

		switch {
		case key.Matches(message, model.keyMap.Menu):
			if !model.menuSpawned {
				model.spawnMenu()
				return model, nil
			}

			model.ActiveTab = model.PreviousTab
			model.menuSpawned = false
			model.tabs = model.tabs[:len(model.tabs)-1]

			if model.ActiveTab == 0 {
				return model, model.tabs[0].model.Init()
			}
			return model, nil
		}
	}

	var command tea.Cmd
	model.tabs[model.ActiveTab].model, command =
		model.tabs[model.ActiveTab].model.Update(message)

	if command != nil {
		tabMessage := command()
		switch tabMessage := tabMessage.(type) {
		case MenuEntry:
			model.switchTab(tabMessage)
		case string:
			switch tabMessage {
			case "switchHome":
				model.switchHome()
			}
		}
	}

	return model, command

}

func (model *Model) switchHome() {
	model.tabs = model.tabs[:len(model.tabs)-1]
	model.menuSpawned = false
	model.ActiveTab = 0
	model.PreviousTab = 0
}

func (model *Model) switchTab(message MenuEntry) {
	model.tabs[model.ActiveTab] = TabEntry{
		title:        message.title,
		model:        message.initialiser(model),
		lastAccessed: time.Now(),
	}
	model.menuSpawned = false
}

// Main function to render contents of the application.
func (model Model) View() string {
	return model.tabs[model.ActiveTab].model.View()
}
