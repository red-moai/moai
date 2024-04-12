package external

type MoaiModel interface {
	ModKey() string
	TerminalWidth() int
	TerminalHeight() int
	AvailableHeight() int
	AvailableWidth() int
	SetTabTitle(string)
	IsReady()bool 
}
