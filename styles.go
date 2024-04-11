package main

import (
	"github.com/charmbracelet/lipgloss"
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{
		Light: "#874BFD",
		Dark:  "#414868",
	}
	activeColor = lipgloss.AdaptiveColor{
		Light: "#874BFD",
		Dark:  "#F7768E",
	}

	inactiveBorderStyle = lipgloss.NewStyle().
				Foreground(highlightColor)

	baseTabStyle = lipgloss.NewStyle().
			Padding(0, 1)

	inactiveTabStyle = baseTabStyle.Copy().
				Border(inactiveTabBorder, true).
				BorderForeground(highlightColor)

	activeTabStyle = baseTabStyle.Copy().
			Border(activeTabBorder, true).
			BorderForeground(activeColor)

	windowStyle = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(2, 0).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder()).
			UnsetBorderTop()
)
