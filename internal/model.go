package internal

import (
	"os/user"
	"strings"

	// "time"

	"github.com/Genekkion/moai/apps/home"
	"github.com/Genekkion/moai/bubblegum/fzf"
	"github.com/Genekkion/moai/colors"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bgList "github.com/genekkion/bubblegum/list"
)

type Model struct {
	// Globals
	Error    error
	modkey   string
	username string

	ActiveTab   int
	PreviousTab int
	tabs        TabEntries
	statusBar   external.MoaiStatusBar

	keyMap      GlobalKeyMap
	onHome      bool
	listSpawned bool
	fzfSpawned  bool

	menu bgList.List
	fzf  fzf.Fzf

	colorScheme  colors.ColorScheme
	windowWidth  int
	windowHeight int
}

func (model Model) WindowHeight() int {
	return model.windowHeight
}

func (model Model) WindowWidth() int {
	return model.windowWidth
}

var (
	items = []bgList.Item{
		ListItem{
			title: "booya",
		},
		ListItem{
			title: "hello world",
		},
		ListItem{
			title: "i like to move it move it",
		},
	}
)

type ListItem struct {
	title string
}

func (item ListItem) Title() string {
	return item.title
}

// Initialises the model to be ran by bubbletea
func InitModel() Model {
	model := Model{
		onHome:      true,
		listSpawned: false,
		fzfSpawned:  false,
		ActiveTab:   0,
		PreviousTab: 0,
		modkey:      getModkey(),
		colorScheme: colors.NewKanagawa(),
		menu:        bgList.InitDefaultBubblegum(items),
		fzf:         fzf.InitFzf(30, 81, []fzf.Item{}),
	}

	currentUser, err := user.Current()
	if err != nil {
		model.username = "User"
	} else {
		model.username = currentUser.Username
	}

	model.keyMap = initGlobalKeyMap(model.modkey)
	model.tabs = TabEntries{
		{
			title: "Home",
			model: home.InitHome(&model).(home.HomeModel),
		},
	}
	model.statusBar = InitStatusBar(model)

	model.fzf.SetTitle("Moai Apps")
	model.fzf.SetBorder(lipgloss.RoundedBorder())

	return model
}

func (model Model) ModKey() string {
	return model.modkey
}

func (model Model) Username() string {
	return model.username
}

func (model Model) ColorScheme() colors.ColorScheme {
	return model.colorScheme
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("M O A I 🗿"),
		model.tabs[0].model.Init(),
	)
}

// Main function to update contents of application.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	commands := []tea.Cmd{}

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		updatedModel, _ := model.tabs[0].model.Update(message)
		model.tabs[0].model = updatedModel.(external.MoaiApp)

		model.windowHeight = message.Height
		model.windowWidth = message.Width

		updatedModel, _ = model.menu.Update(message)
		model.menu = updatedModel.(bgList.List)
		updatedModel, _ = model.fzf.Update(message)
		model.fzf = updatedModel.(fzf.Fzf)

		updatedModel, _ = model.statusBar.Update(message)
		model.statusBar = updatedModel.(StatusBar)

		updatedModel, _ = model.tabs[model.ActiveTab].model.Update(message)
		model.tabs[model.ActiveTab].model = updatedModel.(external.MoaiApp)

		return model, nil

	case tea.KeyMsg:
		// updatedModel, _ := model.statusBar.Update(external.MoaiStatusBarMessage{
		// 	Title:   "",
		// 	Message: message.String(),
		// })
		// model.statusBar = updatedModel.(StatusBar)

		switch {
		case key.Matches(message, model.keyMap.Quit):
			return model, tea.Quit
		case key.Matches(message, model.keyMap.Menu):
			if !model.fzfSpawned {
				model.listSpawned = !model.listSpawned
			}
			return model, nil
		case key.Matches(message, model.keyMap.NewTab):
			if !model.listSpawned {
				model.fzfSpawned = !model.fzfSpawned
				var command tea.Cmd
				if model.fzfSpawned {
					model.fzf.Clear()
					command = model.fzf.Focus()
				} else {
					model.fzf.Blur()
				}
				return model, command
			}
		}
	}

	if model.listSpawned {
		updatedModel, command := model.menu.Update(message)
		model.menu = updatedModel.(bgList.List)
		commands = append(commands, command)
	} else if model.fzfSpawned {
		updatedModel, command := model.fzf.Update(message)
		model.fzf = updatedModel.(fzf.Fzf)
		commands = append(commands, command)
	}
	updatedModel, command := model.tabs[model.ActiveTab].model.Update(message)
	model.tabs[model.ActiveTab].model = updatedModel.(external.MoaiApp)
	commands = append(commands, command)

	// if command != nil {
	// 	tabMessage := command()
	// 	switch tabMessage := tabMessage.(type) {
	// case SetIndexMessage:
	// 	model.setActiveTab(tabMessage.index)
	// case MenuEntry:
	// 	model.switchTab(tabMessage)
	// case string:
	// 	switch tabMessage {
	// 	case "switchHome":
	// 		model.switchHome()
	// 	}
	// }
	// }

	return model, tea.Batch(commands...)

}

// func (model *Model) switchHome() {
// 	model.tabs = model.tabs[:len(model.tabs)-1]
// 	model.menuSpawned = false
// 	model.ActiveTab = 0
// 	model.PreviousTab = 0
// }
//
// func (model *Model) setActiveTab(newIndex int) {
// 	model.tabs = model.tabs[:len(model.tabs)-1]
// 	model.ActiveTab = newIndex
// 	model.menuSpawned = false
// 	model.onHome = false
// }
//
// func (model *Model) switchTab(message MenuEntry) {
// 	model.tabs[model.ActiveTab] = TabEntry{
// 		title:        message.title,
// 		model:        message.initialiser(model),
// 		lastAccessed: time.Now(),
// 	}
// 	model.menuSpawned = false
// }

// Main function to render contents of the application.
func (model Model) View() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(model.tabs[model.ActiveTab].model.View())
	stringBuilder.WriteByte('\n')
	stringBuilder.WriteString(model.statusBar.View())
	if model.listSpawned {
		model.menu.SetView(stringBuilder.String())
		return model.menu.View()
	} else if model.fzfSpawned {
		model.fzf.SetView(stringBuilder.String())
		return model.fzf.View()
	}
	return stringBuilder.String()
}
