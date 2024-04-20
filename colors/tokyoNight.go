package colors

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TokyoNight struct {
	mode Theme

	foreground1 []lipgloss.Color
	foreground2 []lipgloss.Color
	foreground3 []lipgloss.Color
	background1 []lipgloss.Color
	background2 []lipgloss.Color
	background3 []lipgloss.Color
	accent1     []lipgloss.Color
	accent2     []lipgloss.Color
	accent3     []lipgloss.Color

	ColourMap map[interface{}]lipgloss.Color
}

func NewTokyoNight() ColorScheme {
	return TokyoNight{
		mode: DARK_MODE,

		foreground1: []lipgloss.Color{
			"#343B58",
			"#B4F9F8",
		},
		foreground2: []lipgloss.Color{
			"#0F4B6E",
			"#CFC9C2",
		},
		foreground3: []lipgloss.Color{
			"#166775",
			"#7AA2F7",
		},

		background1: []lipgloss.Color{
			"#CFC9C2",
			"#565F89",
		},
		background2: []lipgloss.Color{
			"#9699A3",
			"#414868",
		},
		background3: []lipgloss.Color{
			"#D6D6DB",
			"#24283B",
		},

		accent1: []lipgloss.Color{
			"#8C4351",
			"#F7768E",
		},
		accent2: []lipgloss.Color{
			"#8F5E15",
			"#E0AF68",
		},
		accent3: []lipgloss.Color{
			"#33635C",
			"#1F1F28",
		},
	}
}

func (scheme TokyoNight) Mode() Theme {
	return scheme.mode
}

func (scheme TokyoNight) SetMode(theme Theme) tea.Cmd {
	scheme.setMode(theme)
	return func() tea.Msg {
		return ColorSchemeChange{
			colorScheme: scheme,
		}
	}
}

func (scheme *TokyoNight) setMode(theme Theme) {
	scheme.mode = theme
}

func (scheme TokyoNight) FG1() lipgloss.Color {
	return scheme.foreground1[scheme.mode]
}

func (scheme TokyoNight) FG2() lipgloss.Color {
	return scheme.foreground2[scheme.mode]
}

func (scheme TokyoNight) FG3() lipgloss.Color {
	return scheme.foreground3[scheme.mode]
}

func (scheme TokyoNight) BG1() lipgloss.Color {
	return scheme.background1[scheme.mode]
}

func (scheme TokyoNight) BG2() lipgloss.Color {
	return scheme.background2[scheme.mode]
}

func (scheme TokyoNight) BG3() lipgloss.Color {
	return scheme.background3[scheme.mode]
}

func (scheme TokyoNight) AC1() lipgloss.Color {
	return scheme.accent1[scheme.mode]
}

func (scheme TokyoNight) AC2() lipgloss.Color {
	return scheme.accent2[scheme.mode]
}

func (scheme TokyoNight) AC3() lipgloss.Color {
	return scheme.accent3[scheme.mode]
}

func (scheme TokyoNight) Set(key any, color lipgloss.Color) {
	scheme.set(key, color)
}

func (scheme *TokyoNight) set(key any, color lipgloss.Color) {
	scheme.ColourMap[key] = color
}

func (scheme TokyoNight) Get(key any) lipgloss.Color {
	return scheme.ColourMap[key]
}
