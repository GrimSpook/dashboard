package switcher

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	preStyle = lipgloss.NewStyle().PaddingLeft(0)

	titleS = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)
)

func (m Model) previewView() string {

	wr := lipgloss.JoinVertical(
		lipgloss.Left,
		m.previewHeader(),
		preStyle.Render(m.previewViewport.View()),
		m.previewFooter(),
	)

	return wr
}

func (m Model) previewHeader() string {
	return "Open Workspaces\n"
}

func (m Model) previewContent() string {
	var str strings.Builder
	for _, open := range m.openWorkspaces {

		str.WriteString(open.Title + "\n")

	}
	return str.String()
}

func (m Model) previewFooter() string {
	return ""
}
