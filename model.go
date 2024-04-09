package main

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// Global
	Tabs       []string
	TabContent []string
	ActiveTab  int // 0 index
	Error      error

	// Home state
	HomeChoices  []string
	HomeCursor   int
	HomeSelected map[int]struct{}
	HomeQuote    string

	// Diary state
	diarySearch textinput.Model
}

func (model *Model) initTabs() {
	// Tab labels
	model.Tabs = []string{
		"Home",
		"Notes",
		"Diary",
		"Settings",
	}
	// TODO: To be removed when each tab's content
	// has been completed
	model.TabContent = []string{
		"home stuff",
		"notes",
		"diary",
		"Settings",
	}
}

// Initialises the model to be ran by bubbletea
func InitModel() Model {
	model := Model{}
	model.initTabs()
	model.initHome()
	model.initDiary()

	return model
}

func (model Model) Init() tea.Cmd {
	return nil
}

const modKey = "alt+"

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := message.(type) {
	case tea.KeyMsg:
		keypress := msg.String()

		if strings.HasPrefix(keypress, modKey) {
			commandKey := strings.TrimPrefix(keypress, modKey)

			tabIndex, err := strconv.Atoi(commandKey)
			if err == nil && tabIndex > 0 && tabIndex <= 9 {
				model.ActiveTab = tabIndex - 1
				return model, nil
			}

			switch commandKey {
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

	switch model.ActiveTab {
	case 0:
		return model.updateHome(message)
	case 2:
		return model.updateDiary(message)
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
	switch model.ActiveTab {
	case 0:
		renderData = model.renderHome()
	case 2:
		renderData = model.renderDiary()
	default:
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
