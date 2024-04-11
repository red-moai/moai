package main

import (
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/Genekkion/moai/apps/home"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type Model struct {
	// Globals
	Error error

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

// Use empty string if no change to title
func (model *Model) SwapActiveModel(title string, replacementModel tea.Model) {
	activeTab := *model.ActiveTab
	if title != "" {
		(*model.Tabs)[activeTab] = title
	}
	model.TabModels[activeTab] = replacementModel
}


// Initialises the model to be ran by bubbletea
func InitModel() Model {
	activeTab := 0
	model := &Model{
		ActiveTab: &activeTab,

		Tabs: &[]string{
			"Home",
		},
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
		//log.Warn("Value entered for MODKEY env missing / invalid, using default of \"alt\"")
		MODKEY = "alt"
	}

	MODKEY += "+"

	//var swapModelFunc SwapModelFunc
	model.TabModels = []tea.Model{
		home.InitHome(model),

		//InitBork(&model),
	}
	return *model
}

func (model Model) ModKey() string {
	return MODKEY
}

func (model Model) Init() tea.Cmd {
	return textinput.Blink
	//return nil
}

// Main function to update contents of application.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		keypress := msg.String()

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
					*model.ActiveTab++
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

// Returns the rows then columns
func getTerminalDimensions() (int, int, error) {
	var dimensions struct {
		rows    uint16
		cols    uint16
		xpixels uint16
		ypixels uint16
	}

	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&dimensions)),
	)
	if err != 0 {
		return 0, 0, err
	}

	return int(dimensions.rows), int(dimensions.cols), err
}

// Main function to render contents of the application.
func (model Model) View() string {
	_, width, _ := getTerminalDimensions()

	doc := strings.Builder{}

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
			templen := int(float64(width)*0.95) -
				lipgloss.Width(row)

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

	doc.WriteString(row)

	templen := int(float64(width)*0.95) -
		lipgloss.Width(row)

	for range templen + 1 {
		doc.WriteString(inactiveBorderStyle.Render("─"))
	}
	doc.WriteString(inactiveBorderStyle.Render("┐"))

	doc.WriteString("\n")

	doc.WriteString(
		windowStyle.Width(int(float64(width) * 0.95)).
			Render(model.TabModels[*model.ActiveTab].View()),
	)
	doc.WriteString("\n")

	return zone.Scan(docStyle.Render(doc.String()))
}
