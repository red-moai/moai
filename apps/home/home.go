package home

import (
	"fmt"
	"math/rand"

	"github.com/Genekkion/moai/external"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	homeQuotes = []string{
		"Yesterday is history, tomorrow is a mystery, but today is a gift. That is why it is called the present.",
		"The fellas",
	}
)

type HomeModel struct {
	HomeChoices  []string
	HomeCursor   int
	HomeSelected map[int]struct{}
	HomeQuote    string
}

func InitHome(_ external.MoaiModel) tea.Model {
	model := HomeModel{}
	model.HomeChoices = []string{
		"Oonga boonga",
		"boo ya",
	}
	model.HomeSelected = make(map[int]struct{})
	model.HomeQuote = homeQuotes[rand.Intn(len(homeQuotes))]
	return model
}

func (model HomeModel) Init() tea.Cmd {
	return nil
}

func (model HomeModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch message.String() {

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			//log.Debug("", "before", model.HomeCursor)
			if model.HomeCursor > 0 {
				model.HomeCursor--
			}
			//log.Debug("", "after", model.HomeCursor, "pressed", message.String())

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			//log.Debug("", "before", model.HomeCursor)
			if model.HomeCursor < len(model.HomeChoices)-1 {
				model.HomeCursor++
			}
			//log.Debug("", "after", model.HomeCursor, "pressed", message.String())
		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := model.HomeSelected[model.HomeCursor]
			if ok {
				delete(model.HomeSelected, model.HomeCursor)
			} else {
				model.HomeSelected[model.HomeCursor] = struct{}{}
			}
		}
	}

	return model, nil
}

func (model HomeModel) View() string {

	// Header
	text := model.HomeQuote + "\n"

	// Iterate over our choices
	for i, choice := range model.HomeChoices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if model.HomeCursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := model.HomeSelected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		text += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	return text
}
