package switcher

import (
	"log"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	input      textinput.Model
	list       []Workspace
	SelectedId string
	width      int
	height     int
	styles     *Styles
	viewport   viewport.Model
	ready      bool
}

func New(workspaces []Workspace, width int, height int) Model {

	i := textinput.New()
	i.Placeholder = "Search…"
	i.SetVirtualCursor(true)
	i.Focus()
	i.CharLimit = 156
	i.SetWidth(width)
	// i.Prompt = "❯ "
	i.Prompt = ""

	m := Model{
		input:      i,
		list:       workspaces,
		width:      width,
		SelectedId: "id2",
		height:     height,
		styles:     DefaultStyles(width+1, height),
	}

	return m
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		widthOffset := 12
		heightOffset := 4
		width := (msg.Width - widthOffset)
		height := (msg.Height - verticalMarginHeight - heightOffset) / 2

		m.input.SetWidth(width)

		if !m.ready {

			m.viewport = viewport.New(viewport.WithWidth(width), viewport.WithHeight(height))
			m.viewport.YPosition = headerHeight

			m.viewport.SetContent(m.listView())

			m.ready = true
		} else {
			m.viewport.SetWidth(width + 1)
			m.viewport.SetHeight(height)

			m.width = width + 1
			m.height = height
		}

	case tea.KeyPressMsg:
		switch msg.String() {

		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)

			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var wrapper string

	if !m.ready {

		wrapper = "\n init..."

	} else {

		wrapper = lipgloss.JoinVertical(lipgloss.Left, m.headerView(), m.viewport.View(), m.footerView())

	}

	return m.styles.view.Render(wrapper)
}

func (m Model) headerView() string {

	str := m.styles.input.Render(m.input.View())

	return str
}

func (m Model) footerView() string {
	return ""
}

func (m *Model) listView() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	var s []string

	for _, item := range m.list {

		// avWidth := (lipgloss.Width(item.Title) + lipgloss.Width(item.Branch) + lipgloss.Width(item.Path))
		// sp := (m.width / 2) - avWidth

		str := item.Title
		title := lipgloss.JoinHorizontal(lipgloss.Center,
			item.Title,
			strings.Repeat(" ", 20-lipgloss.Width(item.Title)),
			item.Branch,
			strings.Repeat(" ", 20-lipgloss.Width(item.Branch)),
			strings.ReplaceAll(item.Path, home, "~"),
		)

		if m.SelectedId == item.Id {
			str = m.styles.listItem.Background(m.styles.selectedColor).Render(title)
		} else {
			str = m.styles.listItem.Render(title)
		}

		s = append(s, str)

	}

	wrapper := lipgloss.JoinVertical(lipgloss.Left, s...)

	return m.styles.list.Width(m.width - 11).Render(wrapper)
}
