package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TabEntry struct {
	title        string
	model        tea.Model
	lastAccessed time.Time
}

func (entry TabEntry) timeElapsedString() string {
	timeElapsed := time.Since(entry.lastAccessed).Round(time.Minute)

	minutes := int(timeElapsed.Minutes())
	hours := minutes / 60
	days := hours / 24

	switch {
	case days > 0:
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case hours > 0:
		if hours == 1 {
			return "1 hour ago"

		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	if minutes < 1 {
		return "Just only"
	}
	return fmt.Sprintf("%d minutes ago", minutes)
}

type TabEntries []TabEntry

func (entries TabEntries) Len() int {
	return len(entries)
}

func (entries TabEntries) Swap(i int, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}

func (entries TabEntries) Less(i int, j int) bool {
	if entries[i].lastAccessed.Equal(entries[j].lastAccessed) {
		return entries[i].title < entries[j].title
	}
	return entries[i].lastAccessed.After(entries[j].lastAccessed)
}
