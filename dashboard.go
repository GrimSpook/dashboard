package main

import (
	"dashboard/data"
	"dashboard/switcher"
	"dashboard/tracker"
	"log"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	switcher switcher.Model
	tracker  tracker.Model
	state    data.DashboardState
	width    int
	height   int
}

func initModel() Model {

	l := data.GenerateSections()
	ol, err := data.GetOpenWorkspaces()
	if err != nil {
		log.Fatal(err)
	}

	sw := switcher.New(l, ol, 70, 30)

	// tr := tracker.New()

	return Model{
		state:    data.StateSwitcher,
		switcher: sw,
		// tracker:  tr,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.switcher.SetHeight(msg.Height)
		m.switcher.SetWidth(msg.Width)
		m.width = msg.Width
		m.height = msg.Height

	case data.UpdateStateMsg:
		m.state = msg.State

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+s":
			if m.state == data.StateSwitcher {
				return m, data.SetState(data.StateTracker, m.width, m.height)
			} else {
				return m, data.SetState(data.StateSwitcher, m.width, m.height)
			}
		}
	}

	var cmd tea.Cmd
	switch m.state {
	case data.StateSwitcher:
		var cmd tea.Cmd
		m.switcher, cmd = m.switcher.Update(msg)
		return m, cmd
	case data.StateTracker:
		var cmd tea.Cmd
		m.tracker, cmd = m.tracker.Update(msg)
		return m, cmd
	}
	return m, cmd
}

func (m Model) View() tea.View {

	switch m.state {
	case data.StateHome:

		v := tea.NewView("Home")
		v.AltScreen = true
		return v

	case data.StateSwitcher:

		v := tea.NewView(m.switcher.View())
		v.AltScreen = true
		return v

	case data.StateTracker:

		v := tea.NewView(m.tracker.View())
		v.AltScreen = true
		return v

	default:
		return tea.NewView("dashboard")
	}
}

func (m Model) titleView() string {

	return ""
}
