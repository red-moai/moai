package main

import (
	"os"

	"github.com/Genekkion/moai/apps/home"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	// Globals
	Error  error
	modkey string

	// Tab
	Tabs        *[]string
	TabModels   []tea.Model
	ActiveTab   *int // 0 index
	keyMap      KeyMap
	HomeModel   home.HomeModel
	menuModel   MenuModel
	onHome      bool
	menuSpawned bool
}

var (
	MAX_TABS       = 9
	VALID_MOD_KEYS = []string{
		"alt",
	}
)

func (model *Model) setModkey() {
	modkeyFound := false
	MODKEY := os.Getenv("MODKEY")
	for _, modKey := range VALID_MOD_KEYS {
		if MODKEY == modKey {
			modkeyFound = true
			break
		}
	}
	if !modkeyFound {
		MODKEY = "alt"
	}
	MODKEY += "+"

	model.modkey = MODKEY
}

// Initialises the model to be ran by bubbletea
func InitModel() Model {
	model := Model{
		onHome:      true,
		menuSpawned: false,
	}

	model.setModkey()
	model.HomeModel = home.InitHome(model).(home.HomeModel)
	model.menuModel = InitMenu(model)
	model.keyMap = initKeyMap(model.modkey)

	return model
}

func (model *Model) toggleMenu() {
	model.menuSpawned = !model.menuSpawned
}

func (model Model) ToggleMenu() {
	model.toggleMenu()
}

func (model Model) ModKey() string {
	return model.modkey
}

func (model Model) GetOnHome() bool {
	return model.onHome
}

func (model Model) SetOnHome(onHome bool) {
	model.onHome = onHome
}

func (model *Model) setOnHome(onHome bool) {
	model.SetOnHome(onHome)
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("M O A I 🗿"),
		model.HomeModel.GetSpinner().Tick,
	)
}

var (
	backToHome = true
)

// Main function to update contents of application.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		keypress := message.String()

		switch keypress {
		case "ctrl+c":
			return model, tea.Quit
		}

		switch {
		case key.Matches(message, model.keyMap.Menu):
			model.toggleMenu()
			backToHome = !model.menuSpawned && model.onHome
		}
	}

	var command tea.Cmd
	switch {
	case model.menuSpawned:
		var menuModel tea.Model
		menuModel, command = model.menuModel.Update(message)
		model.menuModel = menuModel.(MenuModel)

	case model.onHome:
		var homeModel tea.Model
		homeModel, command = model.HomeModel.Update(message)
		model.HomeModel = homeModel.(home.HomeModel)
		if backToHome {
			return model, tea.Batch(command, model.HomeModel.GetSpinner().Tick)
		}
		return model, command
	}

	/*
		model.TabModels[*model.ActiveTab], command =
			model.TabModels[*model.ActiveTab].Update(message)
	*/
	return model, command

}

// Main function to render contents of the application.
func (model Model) View() string {
	if model.menuSpawned {
		return model.menuModel.View()
	} else if model.onHome {
		return model.HomeModel.View()
	}
	return "booya"
}
