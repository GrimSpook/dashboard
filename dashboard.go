package main

import (
	"dashboard/sidebar"
	"dashboard/switcher"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	sidebar  sidebar.Model
	switcher switcher.Model
	width    int
	height   int
}

func initModel() Model {

	items := []sidebar.Item{
		{Title: "󰕮", Id: 0},
		{Title: "󰕮", Id: 1},
		{Title: "󰕮", Id: 2},
		{Title: "󰕮", Id: 3},
	}

	// ws := []switcher.Workspace{
	// 	{Title: "nvim", Path: "~\\nvim", Branch: "main", Id: "id1"},
	// 	{Title: "nvim", Path: "~\\nvim", Branch: "main", Id: "id2"},
	// 	{Title: "nvim", Path: "~\\nvim", Branch: "main", Id: "id3"},
	// 	{Title: "nvim", Path: "~\\nvim", Branch: "main", Id: "id4"},
	// }

	l, err := switcher.FdSearch("dev")
	if err != nil {

	}

	s := sidebar.New(items, 0, 0)
	sw := switcher.New(l, 70, 30)

	return Model{
		sidebar:  s,
		switcher: sw,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.sidebar.SetHeight(msg.Height)
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
	m.sidebar, cmd = m.sidebar.Update(msg)
	m.switcher, cmd = m.switcher.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	var v tea.View
	v.AltScreen = true

	str := lipgloss.JoinHorizontal(lipgloss.Top, m.sidebar.View(), m.switcher.View())

	v.SetContent(str)

	return v
}
