package main

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Choices  []string
	Cursor   int
	Selected map[int]struct{}

	Tabs       []string
	TabContent []string
	ActiveTab  int
}

// Initialises the model to be ran by bubbletea
func InitModel() Model {
	return Model{
		Choices: []string{
			"Oonga boonga",
			"boo ya",
		},
		Selected: make(map[int]struct{}),

		Tabs: []string{
			"Home",
			"Notes",
			"Settings",
		},

		TabContent: []string{
			"home stuff",
			"notes",
			"Settings",
		},
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if model.ActiveTab == 0 {
		return model.updateHome(message)
	}

	switch msg := message.(type) {
	case tea.KeyMsg:
		keypress := msg.String()

		if strings.HasPrefix(keypress, "alt+") {
			tabIndex, err := strconv.Atoi(strings.TrimPrefix(keypress, "alt+"))
			if err != nil || tabIndex == 0 ||
				tabIndex > len(model.TabContent) {
				return model, nil
			}

			model.ActiveTab = tabIndex - 1
			return model, nil
		}

		switch keypress {
		case "ctrl+c", "q":
			return model, tea.Quit
		case "right", "l", "tab":
			model.ActiveTab = min(model.ActiveTab+1, len(model.Tabs)-1)
			return model, nil
		case "left", "h", "shift+tab":
			model.ActiveTab = max(model.ActiveTab-1, 0)
			return model, nil
		}
	}

	return model, nil
}

func (model Model) View() string {

	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range model.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(model.Tabs)-1, i == model.ActiveTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	var renderData string
	if model.ActiveTab == 0 {
		renderData = model.renderHome()
	} else {
		renderData = model.TabContent[model.ActiveTab]
	}

	doc.WriteString(
		windowStyle.Width(
			lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize(),
		).Render(
			//model.TabContent[model.ActiveTab],
			renderData,
		),
	)

	return docStyle.Render(doc.String())
}

func (model Model) updateHome(message tea.Msg) (tea.Model, tea.Cmd) {

	switch message := message.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch message.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return model, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if model.Cursor > 0 {
				model.Cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if model.Cursor < len(model.Choices)-1 {
				model.Cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := model.Selected[model.Cursor]
			if ok {
				delete(model.Selected, model.Cursor)
			} else {
				model.Selected[model.Cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return model, nil
}

func (model Model) renderHome() string {

	// The header
	text := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range model.Choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if model.Cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := model.Selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		text += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	return text
}
