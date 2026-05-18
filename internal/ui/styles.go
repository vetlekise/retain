package ui

import lipgloss "charm.land/lipgloss/v2"

var (
	titleStyle     = lipgloss.NewStyle().Bold(true)
	correctStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	incorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

func bodyStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height-1).Padding(1, 2)
}

func footerStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width)
}
