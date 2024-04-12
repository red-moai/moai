package bork

import (
	"net/http"

	"github.com/Genekkion/moai/components"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type BorkModel struct {
	client         http.Client
	borkCategories []string

	ModelList components.DefaultListModel
	MainModel external.MoaiModel
}

var (
	// Depends on modKey
	borkListKeyBindAdd key.Binding

	borkListKeyBindEscape = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "exit"),
	)

	borkCategories = []string{
		"ChatsTodo",
		"PottySense",
		"New category",
	}
)

func InitBork(mainModel external.MoaiModel) tea.Model {
	borkListKeyBindAdd = key.NewBinding(
		key.WithKeys(mainModel.ModKey()+"a"),
		key.WithHelp(mainModel.ModKey()+"a", "add new"),
	)

	model := BorkModel{
		client: *http.DefaultClient,
		ModelList: components.InitDefaultList(
			borkEntries,
			"bork bork",
			mainModel,
			nil,
			nil,
			borkListKeyBindAdd,
			borkListKeyBindEscape,
		),
		borkCategories: borkCategories,
	}
	return model
}

func (model BorkModel) Init() tea.Cmd {
	return nil
}

func (model BorkModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		if model.ModelList.List.FilterState() == list.Filtering {
			break
		}

		if key.Matches(message, borkListKeyBindAdd) {
			borkFormModel := initBorkForm(&model)
			//model.MainModel.SwapActiveModel("", borkFormModel)

			return borkFormModel, nil
		}
	}

	var command tea.Cmd
	model.ModelList, command = model.ModelList.Update(message)
	return model, command
}

func (model BorkModel) View() string {
	return borkListStyle.Render(model.ModelList.View())
}
