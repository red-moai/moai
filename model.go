package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/Genekkion/moai/apps/home"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// Globals
	Error error

	isReady        *bool
	terminalWidth  int
	terminalHeight int

	// Tab
	Tabs      *[]string
	TabModels []tea.Model
	ActiveTab *int // 0 index
}

var (
	MAX_TABS       = 9
	VALID_MOD_KEYS = []string{
		"alt",
	}

	MODKEY = "alt+"
)

func (model Model) tabView() string {
	var renderedTabs []string
	for i, tab := range *model.Tabs {
		var style lipgloss.Style
		isFirst := i == 0
		isLast := i == len(*model.Tabs)-1
		isActive := i == *model.ActiveTab

		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst {
			if isActive {
				border.BottomLeft = "│"
			} else {
				border.BottomLeft = "├"
			}
		} else if isLast {
			if isActive {
				border.BottomRight = "│"
			} else {
				border.BottomRight = "┤"
			}
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(tab))
		if isLast {
			row := lipgloss.JoinHorizontal(
				lipgloss.Top,
				renderedTabs...)
			templen := max(0, model.terminalWidth-
				lipgloss.Width(row))
			if templen > 0 {
				if isActive {
					border.BottomRight = "└"
				} else {
					border.BottomRight = "┴"
				}
			}

			style = style.Border(border)
			renderedTabs[len(renderedTabs)-1] = style.Render(tab)
		}

	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderedTabs...)
	text := strings.Builder{}
	text.WriteString(row)

	templen := max(0, model.terminalWidth-
		lipgloss.Width(row)-1)

	for range templen {
		text.WriteString(inactiveBorderStyle.Render("─"))
	}
	text.WriteString(inactiveBorderStyle.Render("┐"))

	return text.String()
}

// Initialises the model to be ran by bubbletea
func InitModel() Model {
	activeTab := 0
	isReady := false
	model := Model{
		ActiveTab: &activeTab,

		Tabs: &[]string{
			"Home",
		},
		isReady: &isReady,
	}

	isFound := false
	MODKEY = os.Getenv("MODKEY")

	for _, modKey := range VALID_MOD_KEYS {
		if MODKEY == modKey {
			isFound = true
			break
		}
	}
	if !isFound {
		MODKEY = "alt"
	}

	MODKEY += "+"

	model.TabModels = []tea.Model{
		home.InitHome(model),
	}

	return model
}

func (model Model) ModKey() string {
	return MODKEY
}

func (model Model) TerminalHeight() int {
	return model.terminalHeight
}

func (model Model) TerminalWidth() int {
	return model.terminalWidth
}

func (model Model) Style() lipgloss.Style {
	return mainStyle
}

func (model Model) AvailableHeight() int {
	return model.terminalHeight -
		lipgloss.Height(model.tabView()) -
		3 // 2 for padding, 1 for border
}
func (model Model) AvailableWidth() int {
	return model.terminalWidth -
		4 // 2 for padding, 2 for border
}

func (model Model) SetTabTitle(title string) {
	(*model.Tabs)[*model.ActiveTab] = title
}

func (model Model) Init() tea.Cmd {
	return tea.SetWindowTitle("M O A I 🗿")
}

func (model Model) IsReady() bool {
	return *model.isReady
}

// Main function to update contents of application.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:

		mainStyle.Height(message.Height -
			lipgloss.Height(model.tabView()))
		mainStyle.Width(message.Width - 2)

		model.terminalHeight = message.Height
		model.terminalWidth = message.Width

		*model.isReady = true

	case tea.KeyMsg:
		keypress := message.String()

		if strings.HasPrefix(keypress, MODKEY) {
			commandKey := strings.TrimPrefix(keypress, MODKEY)

			tabIndex, err := strconv.Atoi(commandKey)
			if err == nil && tabIndex > 0 && tabIndex <= 9 {
				*model.ActiveTab = tabIndex - 1
				return model, nil
			}

			switch commandKey {
			case "t":
				// Hard cap for number of tabs
				if len(*(model.Tabs)) < MAX_TABS {
					*model.Tabs = append(*model.Tabs, "New Tab")
					model.TabModels = append(
						model.TabModels,
						InitNewTab(&model),
					)

					// Move right
					*model.ActiveTab = len(*(model.Tabs)) - 1
				}

				return model, nil
			case "c":
				if len(*model.Tabs) == 1 {
					return model, tea.Quit
				}

				*model.Tabs = append(
					(*model.Tabs)[:*model.ActiveTab],
					(*model.Tabs)[*model.ActiveTab+1:]...,
				)
				model.TabModels = append(
					model.TabModels[:*model.ActiveTab],
					model.TabModels[*model.ActiveTab+1:]...,
				)

				// Move left
				if *model.ActiveTab > 0 {
					*model.ActiveTab--
				}
				return model, nil

			case "q":
				return model, tea.Quit

			case "right", "l":
				if *model.ActiveTab < len(*model.Tabs)-1 {
					*model.ActiveTab++
				}

				return model, nil

			case "left", "h":
				if *model.ActiveTab > 0 {
					*model.ActiveTab--
				}

				return model, nil
			}

		}

		switch keypress {
		case "ctrl+c":
			return model, tea.Quit
		}

	}

	var command tea.Cmd
	model.TabModels[*model.ActiveTab], command =
		model.TabModels[*model.ActiveTab].Update(message)
	return model, command

}

// Main function to render contents of the application.
func (model Model) View() string {
	text := strings.Builder{}
	//text.WriteString(model.tabView())
	text.WriteString(model.TabModels[*model.ActiveTab].View()
	)
	return text.String()
}
