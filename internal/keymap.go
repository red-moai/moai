package internal

import "github.com/charmbracelet/bubbles/key"

type GlobalKeyMap struct {
	Quit   key.Binding
	Menu   key.Binding
	Home   key.Binding
	NewTab key.Binding
}

func initGlobalKeyMap(modkey string) GlobalKeyMap {
	keymap := GlobalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Menu: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "toggle menu"),
		),
		NewTab: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "new tab"),
		),
	}

	return keymap
}

type MenuKeyMap struct {
	Exit  key.Binding
	Help  key.Binding
	Focus key.Binding
	Home  key.Binding
}

func initMenuKeyMap(modkey string) MenuKeyMap {
	keymap := MenuKeyMap{
		Exit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Focus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch focus"),
		),
		Home: key.NewBinding(
			key.WithKeys(modkey+"h"),
			key.WithHelp(modkey+"h", "home"),
		),
	}

	return keymap
}

func (keymap MenuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keymap.Home, keymap.Focus, keymap.Help, keymap.Exit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (keymap MenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keymap.Help},
		{keymap.Exit},
	}
}
