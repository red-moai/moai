package bork

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
)

var (
	borkEntries = []list.Item{
		BorkEntry{
			id:          0,
			category:    "Check connection",
			title:       "Httpbun GET",
			description: "Server to check connection health",
			method:      "GET",
			url:         "https://httpbun.com/get",
			response:    "- No record found -",
		},
		BorkEntry{
			id:          1,
			category:    "Check connection",
			title:       "Httpbun PUT",
			description: "Server to check connection health",
			method:      "PUT",
			url:         "https://httpbun.com/put",
			response:    "- No record found -",
		},
		BorkEntry{
			id:          2,
			category:    "Check connection",
			title:       "Httpbun POST",
			description: "Server to check connection health",
			method:      "POST",
			url:         "https://httpbun.com/post",
			response:    "- No record found -",
		},
		BorkEntry{
			id:          3,
			category:    "Check connection",
			title:       "Httpbun DELETE",
			description: "Server to check connection health",
			method:      "DELETE",
			url:         "https://httpbun.com/delete",
			response:    "- No record found -",
		},
		BorkEntry{
			id:          4,
			category:    "Check connection",
			title:       "Httpbun PATCH",
			description: "Server to check connection health",
			method:      "PATCH",
			url:         "https://httpbun.com/patch",
			response:    "- No record found -",
		},
	}
)

type BorkEntry struct {
	id           int
	title        string
	method       string
	category     string
	description  string
	url          string
	response     string
	requestStart time.Time
}

func (entry BorkEntry) Title() string {
	return entry.title
}
func (entry BorkEntry) Description() string {
	return entry.description
}
func (entry BorkEntry) FilterValue() string {
	return entry.title
}
