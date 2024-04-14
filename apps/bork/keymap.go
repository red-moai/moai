package bork

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	AddNew      key.Binding
	SendRequest key.Binding
	Exit        key.Binding
}

func initKeyMap(modkey string) KeyMap {
	return KeyMap{
		Exit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit"),
		),
		AddNew: key.NewBinding(
			key.WithKeys(modkey+"a"),
			key.WithHelp(modkey+"a", "add new"),
		),
		SendRequest: key.NewBinding(
			key.WithKeys(modkey+"enter"),
			key.WithHelp(modkey+"enter", "add new"),
		),
	}
}
