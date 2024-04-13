package gpt

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Help   key.Binding
	Quit   key.Binding
	Menu   key.Binding
	NewTab key.Binding
}

func (keyMap KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		keyMap.Help,
		keyMap.Quit,
	}
}

func (keyMap KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keyMap.NewTab, keyMap.Menu},
		{keyMap.Help, keyMap.Quit},
	}
}

func initKeyMap(modkey string) KeyMap {
	keys := KeyMap{
		Help: key.NewBinding(
			key.WithKeys(modkey+"h", "?"),
			key.WithHelp(modkey+"h/?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Menu: key.NewBinding(
			key.WithKeys(modkey+"e"),
			key.WithHelp(modkey+"e", "toggle menu"),
		),
		NewTab: key.NewBinding(
			key.WithKeys(modkey+"t"),
			key.WithHelp(modkey+"t", "new tab"),
		),
	}

	return keys
}
