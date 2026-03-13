package sidebar

import (
	// "fmt"
	// "strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Item struct {
	Title string
	Id    int
}

type Model struct {
	items      []Item
	SelectedId int
	width      int
	height     int
	styles     *Styles
}

func New(items []Item, width int, height int) Model {

	m := Model{
		items:      items,
		SelectedId: 0,
		width:      width,
		height:     height,
		styles:     DefaultStyles(),
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SwitchMenuItem() {
	if m.SelectedId == len(m.items)-1 {
		m.SelectedId = 0
	} else {
		m.SelectedId++
	}
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {

		case "tab":
			m.SwitchMenuItem()
		}
	}

	return m, nil
}

func (m Model) View() string {

	var s []string

	for _, item := range m.items {

		title := item.Title

		if m.SelectedId == item.Id {
			title = m.styles.MenuItem.Foreground(m.styles.SelectedColor).Render(item.Title)
		} else {
			title = m.styles.MenuItem.Render(item.Title)
		}

		s = append(s, title)

	}

	wrapper := lipgloss.JoinVertical(lipgloss.Center, s...)

	return m.styles.view.Height(m.height).Render(wrapper)
}
