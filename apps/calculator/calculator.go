package calculator

import (
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type CalculatorModel struct {
	Textarea textarea.Model
	result   string
}

func InitCalculator(_ external.MoaiModel) tea.Model {
	model := CalculatorModel{
		Textarea: textarea.New(),
	}
	model.Textarea.Placeholder = "Whatcha need help with..."
	model.Textarea.Focus()

	model.Textarea.Prompt = "┃ "
	model.Textarea.ShowLineNumbers = false
	return model
}

func (model CalculatorModel) Init() tea.Cmd {
	return nil
}

func (model CalculatorModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var command tea.Cmd
	model.Textarea, command = model.Textarea.Update(message)

	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "enter":

		}

	}

	return model, command
}

func (model CalculatorModel) View() string {
	text := "Quick math\n"
	text += model.Textarea.View() + "\n"
	return text
}
