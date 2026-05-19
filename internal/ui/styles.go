package ui

import lipgloss "charm.land/lipgloss/v2"

var (
	titleStyle     = lipgloss.NewStyle().Bold(true)
	correctStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // ANSI green
	incorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // ANSI red
)

// bodyStyle and footerStyle are functions (not package-level vars) because
// they need the current terminal dimensions passed in on each render.
func bodyStyle(width, height int) lipgloss.Style {
	// Height(height-1) reserves one line for the footer rendered outside the body box.
	return lipgloss.NewStyle().Width(width).Height(height-1).Padding(1, 2)
}

func footerStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width)
}
