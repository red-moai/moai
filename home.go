package main

import (
	"fmt"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	homeQuotes = []string{
		"Today is the present",
		"The fellas",
	}
)

func (model *Model) initHome() {
	model.HomeChoices = []string{
		"Oonga boonga",
		"boo ya",
	}
	model.HomeSelected = make(map[int]struct{})
	model.HomeQuote = homeQuotes[rand.Intn(len(homeQuotes))]
}

func (model Model) updateHome(message tea.Msg) (tea.Model, tea.Cmd) {

	switch message := message.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch message.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return model, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if model.HomeCursor > 0 {
				model.HomeCursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if model.HomeCursor < len(model.HomeChoices)-1 {
				model.HomeCursor++
			}

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

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return model, nil
}

func (model Model) renderHome() string {

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
