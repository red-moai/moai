package template

import (
	"github.com/Genekkion/moai/external"
	tea "github.com/charmbracelet/bubbletea"
)

var ()

type MyModel struct {
}

func InitTodo(_ external.MoaiModel) tea.Model {
	model := MyModel{}

	return model
}

func (model MyModel) Init() tea.Cmd {
	return nil
}

func (model MyModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model, nil
}

func (model MyModel) View() string {
	return "TODO"
}
