package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit   key.Binding
	Menu   key.Binding
	NewTab key.Binding
}

func initKeyMap(modkey string) KeyMap {
	keys := KeyMap{
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
