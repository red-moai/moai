package bork

import (
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

var (
	HTTPOptions = []huh.Option[string]{
		huh.NewOption("GET", http.MethodGet),
		huh.NewOption("POST", http.MethodPost),
		huh.NewOption("PUT", http.MethodPut),
		huh.NewOption("DELETE", http.MethodDelete),
	}
)

type BorkFormModel struct {
	// For swapping back to parent
	// once done
	parentModel *BorkModel

	form     *huh.Form
	category string
}

func initBorkForm(parentModel *BorkModel) tea.Model {

	categoryList := make([]huh.Option[string], len(parentModel.borkCategories))
	for i, category := range parentModel.borkCategories {
		categoryList[i] = huh.NewOption(category, strings.ToLower(category))
	}

	borkFormKeyMap := huh.NewDefaultKeyMap()
	borkFormKeyMap.Text.Editor.Unbind()

	model := BorkFormModel{
		parentModel: parentModel,
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select category").
				Value(&model.category).
				Options(categoryList...),
		),

		huh.NewGroup(
			huh.NewInput().
				Key("newCategory").
				Title("New Category"),
		).WithHideFunc(func() bool {

			return model.category != "new category"
		}),

		huh.NewGroup(
			huh.NewInput().
				Key("title").
				Title("Give it a title"),

			huh.NewText().
				Key("description").
				Title("How about a description"),
		),

		huh.NewGroup(
			huh.NewSelect[string]().
				Key("method").
				Title("Choose the request method").
				Options(HTTPOptions...),

			huh.NewInput().
				Key("url").
				Title("What URL to bork at"),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Key("completed").
				Title("Bork bork?").
				Affirmative("Bork!").
				Negative("Nay"),
		),
	).WithKeyMap(borkFormKeyMap)

	model.form = form
	return model
}

func (model BorkFormModel) Init() tea.Cmd {
	return model.form.NextField()
}

func (model BorkFormModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	commands := []tea.Cmd{}

	if model.form.Get("completed") != nil {
		model.parentModel.ModelList.List.ResetSelected()
		model.parentModel.ModelList.SelectedItem = nil
		commands = append(commands,
			model.parentModel.ModelList.List.InsertItem(
				0,
				BorkEntry{
					title:       model.form.GetString("title"),
					category:    model.form.GetString("category"),
					description: model.form.GetString("description"),
					method:      model.form.GetString("method"),
					url:         model.form.GetString("url"),
				},
			),
		)
		return model.parentModel, tea.Batch(commands...)
	}

	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.Type {
		case tea.KeyEscape:
			model.parentModel.ModelList.List.ResetSelected()
			model.parentModel.ModelList.SelectedItem = nil
			return model.parentModel, tea.Batch(commands...)
		}
	}

	var command tea.Cmd
	formModel, command := model.form.Update(message)
	form, ok := formModel.(*huh.Form)
	if ok {
		model.form = form
		commands = append(commands, command)
	}

	return model, tea.Batch(commands...)
}
func (model BorkFormModel) View() string {
	text := "Who we borkin at\n"
	text += model.form.View() + "\n"
	return text
}
