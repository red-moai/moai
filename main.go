package main

import (
	"fmt"
	"os"

	"github.com/Genekkion/moai/internal/utils"
	"github.com/joho/godotenv"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	utils.LogError(godotenv.Load(), true)

	program := tea.NewProgram(InitModel())
	_, err := program.Run()
	if err != nil {
		fmt.Println("Something went wrong D:<")
		os.Exit(1)
	}

}
