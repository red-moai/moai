package main

import tea "github.com/charmbracelet/bubbletea"

type MoaiModel interface {
	Init() tea.Cmd
	Update(*Model, tea.Msg) (tea.Model, tea.Cmd)
	View(Model) string
}
