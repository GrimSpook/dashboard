package data

import tea "charm.land/bubbletea/v2"

type errMsg struct{ err error }
type UpdateStateMsg struct {
	Width  int
	Height int
	State  DashboardState
}

type DashboardState int

const (
	StateHome DashboardState = iota
	StateSwitcher
	StateTracker
)

func SetState(newState DashboardState, width, height int) tea.Cmd {
	return func() tea.Msg {
		return UpdateStateMsg{
			Width:  width,
			Height: height,
			State:  newState,
		}
	}

}
