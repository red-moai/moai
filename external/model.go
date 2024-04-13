package external

import tea "github.com/charmbracelet/bubbletea"

type MoaiModel interface {
	ModKey() string
	GetLatestWindowMessage() tea.Msg
}
