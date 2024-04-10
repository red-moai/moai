package main

import (
	"errors"
	"os"

	"github.com/Genekkion/moai/internal/log"
	"github.com/joho/godotenv"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

var requiredEnv = []string{
	"MODKEY",
}

func main() {
	log.FatalWrapper(godotenv.Load())

	for _, env := range requiredEnv {
		_, isSet := os.LookupEnv(env)
		if !isSet {
			log.FatalWrapper(errors.New(
				"Required env: \"" + env + "\" not set. Exiting.",
			))
		}
	}

	zone.NewGlobal()
	defer zone.Close()

	program := tea.NewProgram(InitModel())
	_, err := program.Run()
	log.FatalWrapper(err)

}
