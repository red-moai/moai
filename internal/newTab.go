package internal

import (
	"github.com/Genekkion/moai/internal/components"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelInit func(*Model) tea.Model

var (
	MOAI_APPS = map[string]ModelInit{
		"Home":  InitHome,
		"Diary": InitDiary,
		"Bork":  InitBork,
	}

	AVAILABLE_APPS = []list.Item{
		TabEntry{
			title:       "Home",
			description: "Dashboard",
		},
		TabEntry{
			title:       "Diary",
			description: "Your personal diary",
		},
		TabEntry{
			title:       "Bork",
			description: "A HTTP client for quick testing",
		},
	}
)

type NewTabModel struct {
	ModelList components.DefaultListModel
	MainModel *Model
}

type TabEntry struct {
	title       string
	description string
}

func (tabEntry TabEntry) Title() string {
	return tabEntry.title
}
func (tabEntry TabEntry) Description() string {
	return tabEntry.description
}
func (tabEntry TabEntry) FilterValue() string {
	return tabEntry.title + " " + tabEntry.description
}

func InitNewTab(mainModel *Model) NewTabModel {
	model := NewTabModel{
		ModelList: components.InitDefaultList(
			AVAILABLE_APPS,
			"Available apps",
			30,
			30,
			nil,
		),
		MainModel: mainModel,
	}
	return model
}

func (newTabModel NewTabModel) Init() tea.Cmd {
	return nil
}

func (model NewTabModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var command tea.Cmd
	model.ModelList, command = model.ModelList.Update(message)

	// If selected item, swap with current tab
	if model.ModelList.SelectedItem != nil {
		title := model.ModelList.SelectedItem.(TabEntry).title
		replacementModel := MOAI_APPS[title](model.MainModel)

		model.MainModel.SwapActiveModel(title, replacementModel)
		return replacementModel, command

	}
	return model, command
}


func (model NewTabModel) View() string {
	return model.ModelList.View() + "\n"
}
