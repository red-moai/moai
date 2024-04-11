package todo

import (
	"github.com/Genekkion/moai/components"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	addNewKeyBind key.Binding
)

type TodoModel struct {
	list      components.DefaultListModel
	mainModel external.MoaiModel
}

func InitTodo(mainModel external.MoaiModel) tea.Model {
	model := TodoModel{
		list: components.InitDefaultList(
			fakeTodoData,
			"Todo List",
			30,
			30,
			nil,
			nil,
		),
		mainModel: mainModel,
	}

	addNewKeyBind = key.NewBinding(
		key.WithKeys(mainModel.ModKey()+"a"),
		key.WithHelp(mainModel.ModKey()+"a", "add new"),
	)

	return model
}

func (model TodoModel) Init() tea.Cmd {
	return nil
}

func (model TodoModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		if key.Matches(message, addNewKeyBind) {
			form := InitForm(model)
			return form, nil
		}
	}

	var command tea.Cmd
	model.list, command = model.list.Update(message)
	return model, command
}

func (model TodoModel) View() string {
	return model.list.View()
}
