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

	s.selectedColor = lipgloss.Color("10")

	s.listItem = lipgloss.NewStyle()

	s.view = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).Padding(0, 1).
		BorderForeground(lipgloss.Color("#ff0000")).
		Background(lipgloss.Color("#300000"))

	s.list = lipgloss.NewStyle()

	s.input = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false)

	return s
}
