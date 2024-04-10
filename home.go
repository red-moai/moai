package main

import (
	"fmt"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	homeQuotes = []string{
		"Yesterday is history, tomorrow is a mystery, but today is a gift. That is why it is called the present.",
		"The fellas",
	}
)

type HomeModel struct{}

func (model *Model) initHome() {
	model.HomeChoices = []string{
		"Oonga boonga",
		"boo ya",
	}
	model.HomeSelected = make(map[int]struct{})
	model.HomeQuote = homeQuotes[rand.Intn(len(homeQuotes))]
}

func (model HomeModel) Init() tea.Cmd {
	return nil
}

func (homeModel HomeModel) Update(model *Model, message tea.Msg) (tea.Model, tea.Cmd) {
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

	//log.Debug("", "exiting", model.HomeCursor)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return model, nil
}

func (homeModel HomeModel) View(model Model) string {
	height, width, _ := getTerminalDimensions()

	// Header
	text := model.HomeQuote + "\n"
	text += fmt.Sprintf("height: %d \n", height)
	text += fmt.Sprintf("width: %d \n", width)

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
