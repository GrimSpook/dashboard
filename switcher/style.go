package switcher

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type Styles struct {
	selectedColor color.Color
	listItem      lipgloss.Style
	view          lipgloss.Style
	input         lipgloss.Style
	list          lipgloss.Style
}

func DefaultStyles(width, height int) *Styles {

	s := new(Styles)

	s.selectedColor = lipgloss.Color("#db6363")

	s.listItem = lipgloss.NewStyle()

	s.view = lipgloss.NewStyle()
	// Border(lipgloss.NormalBorder()).Padding(0, 1).
	// BorderForeground(lipgloss.Color("8"))

	s.list = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).Padding(0, 1).
		BorderForeground(lipgloss.Color("8"))

	s.input = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("#ff605a"))

	return s
}
