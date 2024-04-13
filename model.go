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

	ActiveTab   int // 0 index
	PreviousTab int
	keyMap      GlobalKeyMap
	onHome      bool
	menuSpawned bool

	latestWindowMsg tea.Msg

	tabs      []string
	tabModels []tea.Model
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
		PreviousTab: 0,
	}

	model.setModkey()
	model.keyMap = initGlobalKeyMap(model.modkey)
	model.tabs = []string{
		"Home",
	}
	model.tabModels = []tea.Model{
		home.InitHome(&model).(home.HomeModel),
	}

	return model
}

func (model Model) ModKey() string {
	return model.modkey
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("M O A I 🗿"),
		model.tabModels[0].Init(),
	)
}

func (model *Model) spawnMenu() {
	model.tabModels = append(model.tabModels, InitMenu(model))
	model.PreviousTab = model.ActiveTab
	model.ActiveTab = len(model.tabModels) - 1
	model.menuSpawned = true
}

func (model *Model) switchBack() {
	model.ActiveTab = model.PreviousTab
	model.menuSpawned = false
}

// Main function to update contents of application.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.tabModels[0], _ = model.tabModels[0].Update(message)

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
				//model.spawnMenu()
				model.tabModels = append(model.tabModels, InitMenu(&model))
				model.PreviousTab = model.ActiveTab
				model.ActiveTab = len(model.tabModels) - 1
				model.menuSpawned = true

				return model, nil
			}

			//model.switchBack()
			model.ActiveTab = model.PreviousTab
			model.menuSpawned = false

			if model.ActiveTab == 0 {
				return model, model.tabModels[0].Init()
			}
			return model, nil
		}
	}

	var command tea.Cmd
	model.tabModels[model.ActiveTab], command =
		model.tabModels[model.ActiveTab].Update(message)

	if command != nil {
		tabMessage := command()
		if tabMessage, ok := tabMessage.(MenuEntry); ok {
			model.switchTab(tabMessage.title, tabMessage.initialiser)
		}
	}

	return model, command

}

func (model *Model) switchTab(title string, initialiser ModelInit) {
	model.tabs = append(model.tabs, title)
	moaiApp := initialiser(model)
	model.tabModels[model.ActiveTab] = moaiApp
	model.menuSpawned = false
}

// Main function to render contents of the application.
func (model Model) View() string {
	return model.tabModels[model.ActiveTab].View()
}
