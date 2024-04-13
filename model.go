package main

import (
	"os"

	"github.com/Genekkion/moai/apps/home"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	// Globals
	Error  error
	modkey string

	CurrentModel *external.MoaiApp
	PrevModel    *external.MoaiApp
	ActiveTab    int // 0 index
	keyMap       GlobalKeyMap
	homeModel    external.MoaiApp
	menuModel    external.MoaiApp
	onHome       bool
	menuSpawned  bool
}

var (
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
		ActiveTab:   0,
	}

	model.setModkey()
	model.homeModel = home.InitHome(model).(home.HomeModel)
	model.menuModel = InitMenu(&model)
	model.keyMap = initGlobalKeyMap(model.modkey)

	model.CurrentModel = &model.homeModel
	model.PrevModel = nil

	return model
}

func (model *Model) toggleMenu() {
	if model.menuSpawned {
		// If menu is visible, then go back to previous screen
		model.CurrentModel = model.PrevModel
		model.PrevModel = nil
	} else {
		// Else spawn a menu
		model.PrevModel = model.CurrentModel
		model.CurrentModel = &model.menuModel
	}
	model.menuSpawned = !model.menuSpawned
}

func (model Model) ModKey() string {
	return model.modkey
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("M O A I 🗿"),
		model.homeModel.Init(),
	)
}

// Main function to update contents of application.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		var modelPointer tea.Model
		modelPointer, _ = model.menuModel.Update(message)
		model.menuModel = modelPointer.(MenuModel)
		modelPointer, _ = model.homeModel.Update(message)
		model.homeModel = modelPointer.(home.HomeModel)

	case tea.KeyMsg:
		keypress := message.String()

		switch keypress {
		case "ctrl+c":
			return model, tea.Quit
		}

		switch {
		case key.Matches(message, model.keyMap.Menu):
			model.toggleMenu()

			if model.onHome {
				return model, model.homeModel.Init()
			}

			return model, nil
		}
	}

	currentModel, command := (*model.CurrentModel).Update(message)
	moaiApp := currentModel.(external.MoaiApp)
	model.CurrentModel = &moaiApp

	return model, command

}

// Main function to render contents of the application.
func (model Model) View() string {
	return (*model.CurrentModel).View()
}
