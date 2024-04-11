package bork

import (
	"net/http"
	"strings"

	"github.com/Genekkion/moai/components"
	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	borkListStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#CFC9C2"))

	// borkSubStyle = lipgloss.NewStyle().
	// 		Background(lipgloss.Color("#FFFFFF"))

	borkEntries = []list.Item{
		BorkEntry{
			category:    "ChatsTodo",
			title:       "Check health",
			description: "Check health from the ChatsTodo backend servers",
			method:      "GET",
			url:         "www.google.com",
		},
		BorkEntry{
			category:    "ChatsTodo",
			title:       "Get summary",
			description: "Get summary from the ChatsTodo backend servers",
			method:      "GET",
			url:         "www.google.com",
		},
		BorkEntry{
			category:    "ChatsTodo",
			title:       "Update stuff",
			description: "Update data in the ChatsTodo backend servers",
			method:      "POST",
			url:         "www.google.com",
		},
	}

	borkCategories = []string{
		"ChatsTodo",
		"PottySense",
	}
)

type BorkEntry struct {
	title       string
	method      string
	category    string
	description string
	url         string
}

func (borkEntry BorkEntry) Title() string {
	return borkEntry.title
}
func (borkEntry BorkEntry) Description() string {
	return borkEntry.description
}
func (borkEntry BorkEntry) FilterValue() string {
	return borkEntry.title + " " + borkEntry.description
}

type BorkModel struct {
	client         http.Client
	borkCategories []string

	subModel *BorkFormModel

	ModelList components.DefaultListModel
	MainModel external.MoaiModel
}

type BorkFormModel struct {
	// For swapping back to parent
	// once done
	parentModel *BorkModel

	form *huh.Form
}

var (
	borkHTTPOptions = []huh.Option[string]{
		huh.NewOption("GET", http.MethodGet),
		huh.NewOption("POST", http.MethodPost),
		huh.NewOption("PUT", http.MethodPut),
		huh.NewOption("DELETE", http.MethodDelete),
	}

	borkKeyMap = huh.NewDefaultKeyMap()
)

func initSubBork(parentModel *BorkModel) tea.Model {
	model := BorkFormModel{
		parentModel: parentModel,
	}
	categoryList := make([]huh.Option[string], len(parentModel.borkCategories))
	for i, category := range parentModel.borkCategories {
		categoryList[i] = huh.NewOption(category, strings.ToLower(category))
	}

	borkKeyMap.Text.Editor.Unbind()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("category").
				Title("Select category").
				Options(categoryList...),

			huh.NewInput().
				Key("title").
				Title("Give it a title"),

			huh.NewText().
				Key("description").
				Title("How about a description"),

			huh.NewSelect[string]().
				Key("method").
				Title("Choose the request method").
				Options(borkHTTPOptions...),

			huh.NewInput().
				Key("url").
				Title("What URL to bork at"),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Key("completed").
				Title("Bork now?").
				Affirmative("Bork!").
				Negative("Nay"),
		),
	).WithKeyMap(borkKeyMap)

	model.form = form
	return model
}

func (model BorkFormModel) Init() tea.Cmd {
	return model.form.NextField()
}

func (model BorkFormModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	commands := []tea.Cmd{}
	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.Type {
		case tea.KeyEscape:
			model.parentModel.ModelList.List.ResetSelected()
			model.parentModel.ModelList.SelectedItem = nil
			model.parentModel.subModel = nil
			model.parentModel.MainModel.SwapActiveModel("", model.parentModel)
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

	switch model.form.State {
	case huh.StateNormal:
	case huh.StateAborted:
	case huh.StateCompleted:
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
		model.parentModel.subModel = nil
		model.parentModel.MainModel.SwapActiveModel("", model.parentModel)
		return model.parentModel, tea.Batch(commands...)
	}

	return model, tea.Batch(commands...)
}
func (model BorkFormModel) View() string {
	text := "Who we borkin at\n"
	text += model.form.View() + "\n"
	return text
}

var (
	// Depends on modKey
	borkListKeyBindAdd key.Binding

	borkListKeyBindEscape = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "exit"),
	)
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
			30,
			30,
			&borkListStyle,
			borkListKeyBindAdd,
			borkListKeyBindEscape,
		),
		subModel:       nil,
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
		switch {
		case key.Matches(message, borkListKeyBindAdd):
			subBorkModel := initSubBork(&model)
			model.MainModel.SwapActiveModel("", subBorkModel)

			return subBorkModel, nil
		}
	}

	var command tea.Cmd
	model.ModelList, command = model.ModelList.Update(message)

	return model, command
}

func (model BorkModel) View() string {
	text := "Bork.\n"
	text += model.ModelList.View() + "\n"
	return text
}
