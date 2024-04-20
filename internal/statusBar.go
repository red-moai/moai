package internal

import (
	"strings"

	"github.com/Genekkion/moai/external"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	MODE_NORMAL = "N"
	MODE_INSERT = "I"
)

type StatusBar struct {
	mode       string
	title      string
	message    string
	modeStyle  lipgloss.Style
	topStyle   lipgloss.Style
	botStyle   lipgloss.Style
	modeColors map[string]lipgloss.Color
}

func InitStatusBar(mainModel Model) StatusBar {
	colorScheme := mainModel.ColorScheme()

	return StatusBar{
		message: "good morning",
		mode: MODE_NORMAL,
		modeStyle: lipgloss.NewStyle().
			Background(colorScheme.FG2()).
			Foreground(colorScheme.BG3()).
			Padding(0, 1).
			Height(1),
		topStyle: lipgloss.NewStyle().
			Background(colorScheme.AC3()).
			Foreground(colorScheme.FG2()).
			Padding(0, 1).
			Height(1),
		botStyle: lipgloss.NewStyle().
			Background(colorScheme.BG3()).
			Foreground(colorScheme.FG2()).
			Height(1),
		modeColors: map[string]lipgloss.Color{
			MODE_NORMAL: colorScheme.FG2(),
			MODE_INSERT: colorScheme.AC1(),
		},
	}

}

func (model *StatusBar) setMessage(message string) {
	model.message = message
}

func (model StatusBar) SetMessage(message string) {
	model.setMessage(message)
}

func (model *StatusBar) setTitle(title string) {
	model.title = title
}

func (model StatusBar) SetTitle(title string) {
	model.setTitle(title)
}

func (model StatusBar) Init() tea.Cmd {
	return nil
}

func (model StatusBar) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {

	case external.MoaiStatusBarMessage:
		if message.Title != "" {
			model.title = message.Title
		}
		model.message = message.Message

	case tea.WindowSizeMsg:

		if message.Width >= model.modeStyle.GetWidth() {
			model.modeStyle.Width(
				min(3, message.Width),
			)
		} else {
			model.modeStyle.Width(
				min(model.modeStyle.GetWidth(), message.Width),
			)
		}

		model.topStyle = model.topStyle.
			Width(max(message.Width-model.modeStyle.GetWidth(), 0))
		model.botStyle = model.botStyle.
			Width(message.Width)

	}

	return model, nil
}

func (model StatusBar) View() string {
	stringBuilder := strings.Builder{}

	stringBuilder.WriteString(model.modeStyle.Render(model.mode))

	stringBuilder.WriteString(model.topStyle.Render("Moai 🗿"))

	stringBuilder.WriteByte('\n')

	stringBuilder.WriteString(model.botStyle.Render(model.message))

	return stringBuilder.String()
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
