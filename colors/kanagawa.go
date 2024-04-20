package colors

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Kanagawa struct {
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

func NewKanagawa() ColorScheme {
	return Kanagawa{
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
			"#1F1F28",
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
			"#16161D",
		},
	}
}

/*
 sumiInk0 = "#16161D",
    sumiInk1 = "#181820",
    sumiInk2 = "#1a1a22",
    sumiInk3 = "#1F1F28",
    sumiInk4 = "#2A2A37",
    sumiInk5 = "#363646",
    sumiInk6 = "#54546D", --fg

    -- Popup and Floats
    waveBlue1 = "#223249",
    waveBlue2 = "#2D4F67",

    -- Diff and Git
    winterGreen = "#2B3328",
    winterYellow = "#49443C",
    winterRed = "#43242B",
    winterBlue = "#252535",
    autumnGreen = "#76946A",
    autumnRed = "#C34043",
    autumnYellow = "#DCA561",

    -- Diag
    samuraiRed = "#E82424",
    roninYellow = "#FF9E3B",
    waveAqua1 = "#6A9589",
    dragonBlue = "#658594",

    -- Fg and Comments
    oldWhite = "#C8C093",
    fujiWhite = "#DCD7BA",
    fujiGray = "#727169",

    oniViolet = "#957FB8",
    oniViolet2 = "#b8b4d0",
    crystalBlue = "#7E9CD8",
    springViolet1 = "#938AA9",
    springViolet2 = "#9CABCA",
    springBlue = "#7FB4CA",
    lightBlue = "#A3D4D5", -- unused yet
    waveAqua2 = "#7AA89F", -- improve lightness: desaturated greenish Aqua

    -- waveAqua2  = "#68AD99",
    -- waveAqua4  = "#7AA880",
    -- waveAqua5  = "#6CAF95",
    -- waveAqua3  = "#68AD99",

    springGreen = "#98BB6C",
    boatYellow1 = "#938056",
    boatYellow2 = "#C0A36E",
    carpYellow = "#E6C384",

    sakuraPink = "#D27E99",
    waveRed = "#E46876",
    peachRed = "#FF5D62",
    surimiOrange = "#FFA066",
    katanaGray = "#717C7C",
*/

func (scheme Kanagawa) Mode() Theme {
	return scheme.mode
}

func (scheme Kanagawa) SetMode(theme Theme) tea.Cmd {
	scheme.setMode(theme)
	return func() tea.Msg {
		return ColorSchemeChange{
			colorScheme: scheme,
		}
	}
}

func (scheme *Kanagawa) setMode(theme Theme) {
	scheme.mode = theme
}

func (scheme Kanagawa) FG1() lipgloss.Color {
	return scheme.foreground1[scheme.mode]
}

func (scheme Kanagawa) FG2() lipgloss.Color {
	return scheme.foreground2[scheme.mode]
}

func (scheme Kanagawa) FG3() lipgloss.Color {
	return scheme.foreground3[scheme.mode]
}

func (scheme Kanagawa) BG1() lipgloss.Color {
	return scheme.background1[scheme.mode]
}

func (scheme Kanagawa) BG2() lipgloss.Color {
	return scheme.background2[scheme.mode]
}

func (scheme Kanagawa) BG3() lipgloss.Color {
	return scheme.background3[scheme.mode]
}

func (scheme Kanagawa) AC1() lipgloss.Color {
	return scheme.accent1[scheme.mode]
}

func (scheme Kanagawa) AC2() lipgloss.Color {
	return scheme.accent2[scheme.mode]
}

func (scheme Kanagawa) AC3() lipgloss.Color {
	return scheme.accent3[scheme.mode]
}

func (scheme Kanagawa) Set(key any, color lipgloss.Color) {
	scheme.set(key, color)
}

func (scheme *Kanagawa) set(key any, color lipgloss.Color) {
	scheme.ColourMap[key] = color
}

func (scheme Kanagawa) Get(key any) lipgloss.Color {
	return scheme.ColourMap[key]
}
