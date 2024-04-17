package external

import (
	"github.com/Genekkion/moai/colors"
	tea "github.com/charmbracelet/bubbletea"
)

type MoaiModel interface {
	ModKey() string
	GetLatestWindowMessage() tea.Msg
	ColorScheme() colors.ColorScheme
	Username() string
}
