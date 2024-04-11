package bork

import "github.com/charmbracelet/bubbles/list"

var (
	borkEntries = []list.Item{
		BorkEntry{
			category:    "ChatsTodo",
			title:       "Check health",
			description: "Check health from the ChatsTodo backend servers",
			method:      "GET",
			url:         "www.google.com",
		},
		BorkEntry{
			category:    "ChatsTodo",
			title:       "Get summary",
			description: "Get summary from the ChatsTodo backend servers",
			method:      "GET",
			url:         "www.google.com",
		},
		BorkEntry{
			category:    "ChatsTodo",
			title:       "Update stuff",
			description: "Update data in the ChatsTodo backend servers",
			method:      "POST",
			url:         "www.google.com",
		},
	}
)

type BorkEntry struct {
	title       string
	method      string
	category    string
	description string
	url         string
}

func (borkEntry BorkEntry) Title() string {
	return borkEntry.title
}
func (borkEntry BorkEntry) Description() string {
	return borkEntry.description
}
func (borkEntry BorkEntry) FilterValue() string {
	return borkEntry.title
}
