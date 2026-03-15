package main

import (
	"dashboard/switcher"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	// sidebar  sidebar.Model
	switcher switcher.Model
	width    int
	height   int
}

func initModel() Model {

	// items := []sidebar.Item{
	// 	{Title: "󰕮", Id: 0},
	// 	{Title: "󰕮", Id: 1},
	// 	{Title: "󰕮", Id: 2},
	// 	{Title: "󰕮", Id: 3},
	// }

	l := switcher.GenerateSections()

	// s := sidebar.New(items, 0, 0)
	sw := switcher.New(l, 70, 30)

	return Model{
		switcher: sw,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// m.sidebar.SetHeight(msg.Height)
		m.switcher.SetHeight(msg.Height)
		m.switcher.SetWidth(msg.Width)
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		}
	}

	var cmd tea.Cmd
	// m.sidebar, cmd = m.sidebar.Update(msg)
	m.switcher, cmd = m.switcher.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	var v tea.View
	v.AltScreen = true

	s := "\n// SEARCH \n"

	str := lipgloss.JoinVertical(lipgloss.Top, s, m.switcher.View())

	wr := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, str)

	v.SetContent(wr)

	return v
}

func (m Model) titleView() string {

	return ""
}
