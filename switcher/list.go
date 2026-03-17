package switcher

import (
	"dashboard/data"

	"github.com/sahilm/fuzzy"
)

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

func (m Model) Filter(query string, workspaces []data.Workspace) []data.Workspace {
	if query == "" {
		return workspaces
	}

	titles := make([]string, len(workspaces))
	for i, w := range workspaces {
		titles[i] = w.Title
	}

	matches := fuzzy.Find(query, titles)

	wsMap := make(map[string]data.Workspace)
	for _, w := range workspaces {
		wsMap[w.Title] = w
	}

	filtered := make([]data.Workspace, 0, len(matches))
	for _, match := range matches {
		filtered = append(filtered, wsMap[match.Str])
	}

	return filtered
}
