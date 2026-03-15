package switcher

type Section struct {
	Title string
	List  []Workspace
}

type Workspace struct {
	Title  string
	Path   string
	Branch string
	Id     string
}

type weztermCliJson struct {
	Window_id float64 `json:"window_id"`
	Tab_id    float64 `json:"tab_id"`
	Pane_id   float64 `json:"pane_id"`
	Workspace string  `json:"workspace"`
	// Size              map[string]size `json:"size"`
	Title             string  `json:"title"`
	Cwd               string  `json:"cwd"`
	Cursor_x          float64 `json:"cursor_x"`
	Cursor_y          float64 `json:"cursor_y"`
	Cursor_shape      string  `json:"cursor_shape"`
	Cursor_visibility string  `json:"cursor_visibility"`
	Left_col          float64 `json:"left_col"`
	Top_row           float64 `json:"top_row"`
	Tab_title         string  `json:"tab_title"`
	Window_title      string  `json:"window_title"`
	Is_active         bool    `json:"is_active"`
	Is_zoomed         bool    `json:"is_zoomed"`
	Tty_name          string  `json:"tty_name"`
}

func (m *Model) scrollToSelected() {
	lineNumber := 0
	found := false

	for _, section := range m.filtered {
		lineNumber++ // Section header line

		for _, ws := range section.List {
			if ws.Id == m.SelectedId {
				found = true
				break
			}
			lineNumber++
		}

		if found {
			break
		}
	}

	if !found {
		return
	}

	viewportHeight := m.viewport.Height()
	currentOffset := m.viewport.YOffset()

	itemTop := lineNumber - 1
	itemBottom := lineNumber + 4

	if itemTop < currentOffset {
		m.viewport.SetYOffset(itemTop)
	}

	if itemBottom > currentOffset+viewportHeight {
		m.viewport.SetYOffset(itemBottom - viewportHeight)
	}
}

func (m Model) moveUp() string {
	for i, section := range m.filtered {
		for j, ws := range section.List {
			if ws.Id != m.SelectedId {
				continue
			}
			if j > 0 {
				return section.List[j-1].Id
			}
			for si := i - 1; si >= 0; si-- {
				prev := m.filtered[si]
				if len(prev.List) == 0 {
					continue
				}
				return prev.List[len(prev.List)-1].Id
			}
		}
	}

	return ""
}

func (m Model) moveDown() string {
	for i, section := range m.filtered {
		for j, ws := range section.List {
			if ws.Id != m.SelectedId {
				continue
			}
			if j < len(section.List)-1 {
				return section.List[j+1].Id
			}
			for si := i + 1; si < len(m.filtered); si++ {
				next := m.filtered[si]
				if len(next.List) == 0 {
					continue
				}
				return next.List[0].Id
			}
		}
	}

	return ""
}
