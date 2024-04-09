package main

import (
	"github.com/Genekkion/moai/internal/utils"
	"github.com/joho/godotenv"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	utils.LogError(godotenv.Load())

	program := tea.NewProgram(InitModel())
	_, err := program.Run()
	utils.LogError(err)

}
