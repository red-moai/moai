package bork

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	AddNew key.Binding
	Exit   key.Binding
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
	}
}
