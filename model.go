package main

import (
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/Genekkion/moai/internal/log"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type Model struct {
	ModKey string

	// Global
	Tabs       []string
	TabContent []string
	ActiveTab  int // 0 index
	TabModels  []MoaiModel
	Error      error

	// New Tabs
	newTabSearch  textinput.Model
	newTabDisplay string

	// Diary state
	diarySearch  textinput.Model
	diaryDisplay string

	// Home state
	HomeChoices  []string
	HomeCursor   int
	HomeSelected map[int]struct{}
	HomeQuote    string
}

var (
	VALID_MOD_KEYS = []string{
		"alt",
	}
)

// Initialises the model to be ran by bubbletea
func InitModel() Model {
	model := Model{}
	model.initTabs()
	model.initHome()
	model.initDiary()

	isFound := false
	model.ModKey = os.Getenv("MODKEY")

	for _, modKey := range VALID_MOD_KEYS {
		if model.ModKey == modKey {
			isFound = true
			break
		}
	}
	if !isFound {
		log.Warn("Invalid MODKEY env entered, using default of \"alt\"")
		model.ModKey = "alt"
	}
	model.ModKey += "+"
	return model
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

		if strings.HasPrefix(keypress, model.ModKey) {
			commandKey := strings.TrimPrefix(keypress, model.ModKey)

			tabIndex, err := strconv.Atoi(commandKey)
			if err == nil && tabIndex > 0 && tabIndex <= 9 {
				model.ActiveTab = tabIndex - 1
				return model, nil
			}

			switch commandKey {
			case "t":
				if len(model.Tabs) < 9 {
					model.Tabs = append(model.Tabs, "New Tab")
					model.TabModels = append(model.TabModels, NewTabModel{})
				}
				return model, nil
			case "c":
				if len(model.Tabs) > 1 {
					model.Tabs = model.Tabs[:len(model.Tabs)-1]
					return model, nil
				}
				return model, tea.Quit
			case "q":
				return model, tea.Quit
			case "right", "l":
				model.ActiveTab = min(model.ActiveTab+1, len(model.Tabs)-1)
				return model, nil
			case "left", "h":
				model.ActiveTab = max(model.ActiveTab-1, 0)
				return model, nil
			}

		}

		switch keypress {
		case "ctrl+c":
			return model, tea.Quit
		}
	}
	return model.TabModels[model.ActiveTab].Update(&model, message)
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
	for i, tab := range model.Tabs {
		var style lipgloss.Style
		isFirst := i == 0
		isLast := i == len(model.Tabs)-1
		isActive := i == model.ActiveTab

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

	if len(model.TabModels) > 0 {
		doc.WriteString(
			windowStyle.Width(int(float64(width) * 0.95)).
				Render(model.TabModels[model.ActiveTab].View(model)),
		)
	} else {
		doc.WriteString(
			windowStyle.Width(int(float64(width) * 0.95)).
				Render(model.TabContent[model.ActiveTab]),
		)
	}

	return zone.Scan(docStyle.Render(doc.String()))
}
