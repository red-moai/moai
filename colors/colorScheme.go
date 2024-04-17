package colors

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Theme int

const (
	LIGHT_MODE = 0
	DARK_MODE  = 1
)

type ColorScheme interface {
	// Get the mode the theme is in
	Mode() Theme
	// Set the mode of the theme. Returns a tea.Cmd
	// containing the ColorSchemeChange struct which contains
	// the updated color scheme model
	SetMode(Theme) tea.Cmd

	// Returns a foreground color
	FG1() lipgloss.Color
	// Returns a foreground color
	FG2() lipgloss.Color
	// Returns a foreground color
	FG3() lipgloss.Color

	// Returns a background color
	BG1() lipgloss.Color
	// Returns a background color
	BG2() lipgloss.Color
	// Returns a background color
	BG3() lipgloss.Color

	// Returns an accent color
	AC1() lipgloss.Color
	// Returns an accent color
	AC2() lipgloss.Color
	// Returns an accent color
	AC3() lipgloss.Color

	// Sets a color in the color map
	Set(any, lipgloss.Color)
	// Get a map from the color map
	Get(any) lipgloss.Color
}

type ColorSchemeChange struct{
	colorScheme ColorScheme
}
