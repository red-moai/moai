package external

import tea "github.com/charmbracelet/bubbletea"

type MoaiModel interface {
	SwapActiveModel(string, tea.Model)
	ModKey() string
}
