package sidebar

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type Styles struct {
	SelectedColor color.Color
	MenuItem      lipgloss.Style
	view          lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)

	s.SelectedColor = lipgloss.Color("1")

	s.MenuItem = lipgloss.NewStyle().PaddingBottom(1)

	s.view = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2)

	return s
}
