package external

import (
	"github.com/Genekkion/moai/colors"
	tea "github.com/charmbracelet/bubbletea"
)

type MoaiModel interface {
	ModKey() string
	WindowWidth() int
	WindowHeight() int
	ColorScheme() colors.ColorScheme
	Username() string
}

type MoaiApp interface {
	tea.Model
}

type MoaiStatusBar interface {
	tea.Model
	SetTitle(string)
	SetMessage(string)
}

type MoaiStatusBarMessage struct {
	Title   string
	Message string
}
