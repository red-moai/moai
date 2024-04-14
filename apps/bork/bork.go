package bork

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Genekkion/moai/external"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BorkModel struct {
	borkCategories []string
	keymap         KeyMap

	list          list.Model
	form          tea.Model
	formShown     bool
	selectedEntry BorkEntry
	selectedIndex int
	hasSelected   bool
	response      string
}

var (
	borkCategories = []string{
		"Check connection",
	}

	modelStyle = lipgloss.NewStyle().
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#CFC9C2"))
	requestStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Border(lipgloss.RoundedBorder())
	requestInfoStyle = lipgloss.NewStyle().
				Padding(0, 2).
				BorderStyle(lipgloss.ThickBorder()).
				BorderBottom(true).
				AlignHorizontal(lipgloss.Left)
	requestHeaderStyle = lipgloss.NewStyle().
				AlignHorizontal(lipgloss.Center).
		//Background(lipgloss.Color("#00FF00")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Foreground(lipgloss.Color("#965027")).
		Bold(true)

	listStyle = lipgloss.NewStyle().
		//Border(lipgloss.NormalBorder())
		Padding(0, 1)
)

func InitBork(mainModel external.MoaiModel) tea.Model {
	model := BorkModel{
		list: list.New(
			borkEntries,
			list.NewDefaultDelegate(),
			30,
			30,
		),
		keymap:         initKeyMap(mainModel.ModKey()),
		borkCategories: borkCategories,
		formShown:      false,
		hasSelected:    false,
	}
	model.list.Title = "Bork Bork! 🐶"

	model.updateDimensions(mainModel.GetLatestWindowMessage())
	model.list.SetShowHelp(false)
	model.list.DisableQuitKeybindings()

	return model
}

func (model BorkModel) Init() tea.Cmd {
	return nil
}

func (model *BorkModel) updateDimensions(message tea.Msg) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		newHeight := message.Height - 2
		newWidth := message.Width - 2
		modelStyle = modelStyle.
			Height(newHeight).
			Width(newWidth)
		newHeight -= 4
		newWidth = (newWidth - 2) / 2

		requestStyle = requestStyle.
			Height(newHeight).
			Width(newWidth)
		requestInfoStyle = requestInfoStyle.
			Width(newWidth)

		listStyle = listStyle.
			Height(newHeight).
			Width(newWidth - 4)
		model.list.SetHeight(newHeight)
		model.list.SetWidth(newWidth - 4)
	}
}

/*
	menuModelStyle = menuModelStyle.
		Height(message.Height - 2).
		Width(message.Width - 2)
	menuHelpStyle = menuHelpStyle.
		Width(message.Width - 4)

	widgetHeight := message.Height - 6 -
		lipgloss.Height(model.helpView())
	widgetWidth := (message.Width-4)/2 - 4

	model.list.SetHeight(widgetHeight)
	model.list.SetWidth(widgetWidth)
	model.table.SetHeight(widgetHeight - 2)

	newColumns := make([]table.Column, len(model.tableColumns))

	for i := range len(model.tableColumns) {
		newColumns[i].Title = model.tableColumns[i].Title
		newColumns[i].Width = ((message.Width-4)/2 - 7) / 2
	}
	model.table.SetColumns(newColumns)

	menuWidgetStyle = menuWidgetStyle.
		Height(widgetHeight - 14).
		Width(widgetWidth + 1)
*/

func formatDuration(duration time.Duration) string {
	if duration < time.Millisecond {
		return fmt.Sprintf("%d ns", duration.Nanoseconds())
	}
	if duration < time.Second {
		return fmt.Sprintf("%.2f ms", float64(duration.Nanoseconds())/float64(time.Millisecond))
	}
	if duration < time.Minute {
		return fmt.Sprintf("%.2f s", float64(duration.Nanoseconds())/float64(time.Second))
	}
	if duration < time.Hour {
		mins := duration / time.Minute
		secs := duration % time.Minute / time.Second
		return fmt.Sprintf("%dmin %ds", mins, secs)
	}
	hours := duration / time.Hour
	mins := duration % time.Hour / time.Minute
	return fmt.Sprintf("%dh %dmin", hours, mins)
}

func (model BorkModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var command tea.Cmd
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.updateDimensions(message)
		return model, nil
	case BorkResponse:
		cursor := model.list.Cursor()
		var updateEntry BorkEntry
		var updateIndex int
		text := fmt.Sprintf(
			"Returned in: %s\n\n%s\n",
			formatDuration(message.timeElapsed),
			message.response,
		)

		for i, entry := range model.list.Items() {
			borkEntry := entry.(BorkEntry)
			if borkEntry.id == message.id {
				borkEntry.response = text
				updateEntry = borkEntry
				updateIndex = i
				break
			}
		}
		command = model.list.SetItem(updateIndex, updateEntry)
		model.list, command = model.list.Update(command)
		model.list.Select(cursor)
		if cursor == model.selectedIndex {
			model.selectedEntry.response = text
		}
		return model, command
	}

	if model.formShown {
		model.form, command = model.form.Update(message)
		if command != nil {
			formMessage := command()
			switch listMessage := formMessage.(type) {
			case BorkEntry:
				model.list.ResetSelected()
				model.list.InsertItem(0, listMessage)
				model.formShown = false
				return model, nil
			}
		}
	} else {
		model.list, command = model.list.Update(message)
		switch message := message.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(message, model.keymap.AddNew):
				model.form = initBorkForm(model.borkCategories)
				model.formShown = true
				return model, nil
			case model.hasSelected && key.Matches(message, model.keymap.SendRequest):
				borkEntry := model.list.SelectedItem().(BorkEntry)
				borkEntry.response = "getting response..."
				borkEntry.requestStart = time.Now()

				command = model.list.SetItem(
					model.list.Cursor(),
					borkEntry,
				)
				model.list, command = model.list.Update(command)
				if model.selectedIndex == model.list.Cursor() {
					model.selectedEntry.response = "getting response..."
					model.selectedEntry.requestStart = time.Now()
				}
				return model, tea.Batch(BorkRequest(model.selectedEntry), command)
			}

			switch message.String() {
			case "enter":
				model.hasSelected = true
				model.selectedIndex = model.list.Cursor()
				model.selectedEntry = model.list.SelectedItem().(BorkEntry)
			}
		}
	}

	return model, command
}

type BorkResponse struct {
	id          int
	response    string
	timeElapsed time.Duration
}

func BorkRequest(entry BorkEntry) tea.Cmd {
	return func() tea.Msg {
		request, err := http.NewRequest(
			entry.method,
			entry.url,
			nil,
		)
		if err != nil {
			return BorkResponse{
				id: entry.id,
				response: fmt.Sprintf(
					"⚠️ Error generating request! ⚠️\n%s\n",
					err.Error(),
				),
				timeElapsed: time.Since(entry.requestStart),
			}
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return BorkResponse{
				id: entry.id,
				response: fmt.Sprintf(
					"⚠️ Error performing request! ⚠️\n%s\n",
					err.Error(),
				),
				timeElapsed: time.Since(entry.requestStart),
			}
		}
		defer response.Body.Close()

		contentType := response.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") {
			var responseData interface{}
			err := json.NewDecoder(response.Body).Decode(&responseData)
			if err != nil {
				return BorkResponse{
					id: entry.id,
					response: fmt.Sprintf(
						"⚠️ Error parsing JSON! ⚠️\n%s\n",
						err.Error(),
					),
					timeElapsed: time.Since(entry.requestStart),
				}
			}
			prettyJson, err := json.MarshalIndent(responseData, "", " ")
			if err != nil {
				return BorkResponse{
					id: entry.id,
					response: fmt.Sprintf(
						"⚠️ Error prettifying JSON! ⚠️\n%s\n",
						err.Error(),
					),
					timeElapsed: time.Since(entry.requestStart),
				}
			}
			return BorkResponse{
				id:          entry.id,
				response:    string(prettyJson),
				timeElapsed: time.Since(entry.requestStart),
			}
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return BorkResponse{
				id: entry.id,
				response: fmt.Sprintf(
					"⚠️ Error reading response string! ⚠️\n%s\n",
					err.Error(),
				),
				timeElapsed: time.Since(entry.requestStart),
			}
		}
		return BorkResponse{
			id:          entry.id,
			response:    string(responseData),
			timeElapsed: time.Since(entry.requestStart),
		}
	}
}

func (model BorkModel) requestView() string {
	text := strings.Builder{}
	text.WriteString(requestHeaderStyle.Render("Request") + "\n")

	if !model.hasSelected {
		text.WriteString("Select an entry to view it here 🗣")
	} else {
		requestInfo := strings.Builder{}
		requestInfo.WriteString(fmt.Sprintf("Category   : %s \n", model.selectedEntry.category))
		requestInfo.WriteString(fmt.Sprintf("Title      : %s \n", model.selectedEntry.title))
		requestInfo.WriteString(fmt.Sprintf("Description: %s \n", model.selectedEntry.description))
		requestInfo.WriteString(fmt.Sprintf("Method     : %s \n", model.selectedEntry.method))
		requestInfo.WriteString(fmt.Sprintf("URL        : %s \n", model.selectedEntry.url))
		text.WriteString(requestInfoStyle.Render(requestInfo.String()))
		text.WriteString("\n")
		tempStyle := lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left)

		text.WriteString(requestHeaderStyle.Render("Response") + "\n")

		text.WriteString(
			tempStyle.Render(
				fmt.Sprintf("%s\n", model.selectedEntry.response),
			),
		)

	}
	return requestStyle.Render(text.String())
}

func (model BorkModel) View() string {
	if model.formShown {
		return modelStyle.Render(model.form.View())
	}

	return modelStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			listStyle.Render(model.list.View()),
			model.requestView(),
		),
	)
}
