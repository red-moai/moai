package main

import "github.com/charmbracelet/bubbles/key"

type GlobalKeyMap struct {
	Quit   key.Binding
	Menu   key.Binding
	NewTab key.Binding
}

func initGlobalKeyMap(modkey string) GlobalKeyMap {
	keys := GlobalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Menu: key.NewBinding(
			key.WithKeys(modkey+"e"),
			key.WithHelp(modkey+"e", "toggle menu"),
		),
		/*
			NewTab: key.NewBinding(
				key.WithKeys(modkey+"t"),
				key.WithHelp(modkey+"t", "new tab"),
			),
		*/
	}

	return keys
}

type MenuKeyMap struct {
	Exit  key.Binding
	Help  key.Binding
	Focus key.Binding
}

func initMenuKeyMap(_ string) MenuKeyMap {
	keys := MenuKeyMap{
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
	}

	return keys
}

func (keymap MenuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keymap.Focus, keymap.Help, keymap.Exit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (keymap MenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keymap.Help},
		{keymap.Exit},
	}
}
