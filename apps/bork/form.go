package bork

import (
	"errors"
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

const (
	NEW_CATEGORY = "New category"
)

type BorkFormModel struct {
	form     *huh.Form
	category string
}

func find(borkCategories []string, target string) bool {
	for _, category := range borkCategories {
		if category == target {
			return true
		}
	}
	return false
}

func initBorkForm(borkCategories []string) tea.Model {

	categoryList := make([]huh.Option[string], len(borkCategories)+1)
	for i, category := range borkCategories {
		categoryList[i] = huh.NewOption(category, strings.ToLower(category))
	}

	categoryList[len(categoryList)-1] = huh.NewOption(
		NEW_CATEGORY,
		NEW_CATEGORY,
	)

	formKeymap := huh.NewDefaultKeyMap()
	formKeymap.Text.Editor.Unbind()
	formKeymap.Quit.Unbind()

	model := BorkFormModel{}

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
				Title("New Category").
				Validate(func(category string) error {
					if category == "" {
						return errors.New("Cannot be left empty!")
					} else if strings.ToLower(category) == strings.ToLower(NEW_CATEGORY) {
						return errors.New("Wait what?")
					} else if find(borkCategories, category) {
						return errors.New("Already exists!")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return model.category != NEW_CATEGORY
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
	).WithKeyMap(formKeymap).
		WithShowHelp(false)

	model.form = form
	return model
}

func (model BorkFormModel) Init() tea.Cmd {
	return model.form.NextField()
}

func (model BorkFormModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	commands := []tea.Cmd{}

	if model.form.Get("completed") != nil {
		// model.parentModel.ModelList.List.ResetSelected()
		// model.parentModel.ModelList.SelectedItem = nil
		// commands = append(commands,
		// 	model.parentModel.ModelList.List.InsertItem(
		// 		0,
		// 		BorkEntry{
		// 		},
		// 	),
		// )
		//return model.parentModel, tea.Batch(commands...)
		return model, func() tea.Msg {
			return BorkEntry{
				title:       model.form.GetString("title"),
				category:    model.form.GetString("category"),
				description: model.form.GetString("description"),
				method:      model.form.GetString("method"),
				url:         model.form.GetString("url"),
			}
		}
	}

	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.Type {
		case tea.KeyEscape:
			// model.parentModel.ModelList.List.ResetSelected()
			// model.parentModel.ModelList.SelectedItem = nil
			//return model.parentModel, tea.Batch(commands...)
			return model, nil
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
