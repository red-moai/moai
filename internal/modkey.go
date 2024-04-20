package internal

import "os"

var (
	VALID_MOD_KEYS = []string{
		"alt",
	}
	MODKEY = "alt+"
)

func getModkey() string {
	modkeyFound := false
	MODKEY := os.Getenv("MODKEY")
	for _, modKey := range VALID_MOD_KEYS {
		if MODKEY == modKey {
			modkeyFound = true
			break
		}
	}
	if !modkeyFound {
		return MODKEY
	}
	MODKEY += "+"
	return MODKEY
}
