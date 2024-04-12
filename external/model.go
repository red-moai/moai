package external

type MoaiModel interface {
	ModKey() string
	GetOnHome() bool
	SetOnHome(bool)
}
