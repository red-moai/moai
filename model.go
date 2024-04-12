package main

import (
	"os"

	"github.com/Genekkion/moai/apps/home"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	// Globals
	Error error

	// Tab
	Tabs      *[]string
	TabModels []tea.Model
	ActiveTab *int // 0 index

	HomeModel home.HomeModel
	onHome    bool
}

var (
	MAX_TABS       = 9
	VALID_MOD_KEYS = []string{
		"alt",
	}

	MODKEY = "alt+"
)

// Initialises the model to be ran by bubbletea
func InitModel() Model {
	model := Model{
		onHome: true,
	}

	modkeyFound := false
	MODKEY = os.Getenv("MODKEY")

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

	model.HomeModel = home.InitHome(model).(home.HomeModel)

	return model
}

func (model Model) ModKey() string {
	return MODKEY
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

// Main function to update contents of application.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		keypress := message.String()

		switch keypress {
		case "ctrl+c":
			return model, tea.Quit
		}
	}

	var command tea.Cmd
	if model.onHome {
		var homeModel tea.Model
		homeModel, command = model.HomeModel.Update(message)
		model.HomeModel = homeModel.(home.HomeModel)
	}

	/*
		model.TabModels[*model.ActiveTab], command =
			model.TabModels[*model.ActiveTab].Update(message)
	*/
	return model, command

}

// Main function to render contents of the application.
func (model Model) View() string {
	//return model.TabModels[*model.ActiveTab].View()
	if model.onHome {
		return model.HomeModel.View()
	}
	return "booya"
}
