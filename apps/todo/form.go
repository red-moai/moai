package todo

import (
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

var ()

type TodoForm struct {
	form        *huh.Form
	parentModel TodoModel
	hasDeadline bool
}

func InitForm(parentModel TodoModel) tea.Model {
	model := TodoForm{
		parentModel: parentModel,
	}

	firstInput := huh.NewInput().
		Key("title").
		Title("Whatcha gotta do?")
	firstInput.Focus()

	form := huh.NewForm(
		huh.NewGroup(
			firstInput,
			huh.NewText().
				Key("description").
				Title("What's it all about?"),

			huh.NewConfirm().
				Title("Got a deadline?").
				Value(&model.hasDeadline).
				Affirmative("Yea...").
				Negative("Nope"),
		),
		huh.NewGroup(
			huh.NewInput().
				Key("deadline").
				Title("When do we gotta do it by?").
				Validate(func(input string) error {
					_, err := time.Parse("02-01-2006 15:04", input)
					if err != nil {
						return errors.New("Dates gotta be like: DD-MM-YYYY HH:MM")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return !model.hasDeadline
		}),

		huh.NewGroup(
			huh.NewConfirm().
				Key("completed").
				Title("Gonna do it?").
				Affirmative("Hell yea").
				Negative("Nuh uh"),
		),
	)
	model.form = form

	return &model
}

func (model TodoForm) Init() tea.Cmd {
	return model.form.Init()
}

func (model TodoForm) Update(message tea.Msg) (tea.Model, tea.Cmd) {

	if model.form.Get("completed") != nil {
		deadline, _ := time.Parse("02-01-2006 15:04", model.form.GetString("deadline"))

		command := model.parentModel.list.List.InsertItem(0,
			TodoEntry{
				title:       model.form.GetString("title"),
				description: model.form.GetString("description"),
				deadline:    deadline,
			},
		)

		return model.parentModel, command
	}

	var command tea.Cmd
	updatedForm, command := model.form.Update(message)
	model.form = updatedForm.(*huh.Form)
	return model, command
}

func (model TodoForm) View() string {
	return model.form.View()
}
