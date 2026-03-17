package switcher

import (
	"dashboard/data"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	input      textinput.Model
	list       []data.Section
	filtered   []data.Section
	SelectedId string
	width      int
	height     int
	styles     *Styles
	viewport   viewport.Model
	ready      bool
}

var (
	borderColorLight = lipgloss.Color("#9f4d58")
	borderColorDark  = lipgloss.Color("#5a313b")

	// borderColorLight = lipgloss.Color("#707070")
	// borderColorDark  = lipgloss.Color("#404040")

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(borderColorLight).
				PaddingRight(1)

	mutedStyle = lipgloss.NewStyle().Foreground(borderColorDark)

	mutedTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	titleStyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).BorderForeground(borderColorLight)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b).Foreground(lipgloss.Color("7"))
	}()

	inline = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		b.Right = "├"
		b.Left = "┤"
		return titleStyle.BorderStyle(b).Foreground(lipgloss.Color("7"))
	}()

	lineStyle = lipgloss.NewStyle().Foreground(borderColorDark)

	pathStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func New(sections []data.Section, width int, height int) Model {

	i := textinput.New()
	i.Placeholder = ""
	i.SetVirtualCursor(true)
	i.Focus()
	i.CharLimit = 156
	i.SetWidth(20)
	// i.Prompt = "❯ "
	i.Prompt = ""

	m := Model{
		input:      i,
		list:       sections,
		filtered:   sections,
		width:      width,
		SelectedId: sections[0].List[0].Id,
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

		widthOffset := 2
		heightOffset := 2
		width := (msg.Width - widthOffset) / 2
		height := (msg.Height - verticalMarginHeight - heightOffset)

		m.input.SetWidth(width - 3)

		if !m.ready {

			m.viewport = viewport.New(viewport.WithWidth(width), viewport.WithHeight(height))
			m.viewport.YPosition = headerHeight

			m.viewport.SetContent(m.listView())

			m.ready = true
		} else {
			m.viewport.SetWidth(width)
			m.viewport.SetHeight(height)

			m.viewport.SetContent(m.listView())

			m.width = width
			m.height = height
		}

	case tea.KeyPressMsg:
		switch msg.String() {

		case "up":
			newId := m.moveUp()
			if newId != "" {
				m.SelectedId = newId
				m.viewport.SetContent(m.listView())
				m.scrollToSelected()
			}

		case "down":
			newId := m.moveDown()
			if newId != "" {
				m.SelectedId = newId
				m.viewport.SetContent(m.listView())
				m.scrollToSelected()
			}

		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)

			value := m.input.Value()
			newSections := []data.Section{}
			for _, section := range m.list {
				section.List = m.Filter(value, section.List)
				if len(section.List) > 0 {
					newSections = append(newSections, section)
				}
			}
			m.filtered = newSections

			if len(m.filtered) > 0 && len(m.filtered[0].List) > 0 {
				m.SelectedId = m.filtered[0].List[0].Id
			}

			m.viewport.SetContent(m.listView())

			return m, cmd
		}
	}

	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	var wrapper string

	if !m.ready {

		wrapper = "\n init..."

	} else {

		wrapper = lipgloss.JoinVertical(lipgloss.Left, m.headerView(), m.viewport.View(), m.footerView())

	}

	b := Border(
		wrapper,
		withWidth(m.viewport.Width()),
		withHeight(m.viewport.Height()),
		withTitleColor(m.styles.selectedColor),
		withCornorColor(borderColorDark),
		withSideColor(borderColorDark),
		withCornorChars(BorderCornorRound),
		withSideChars(BorderSideThin),
		withTitle("Switcher"),
		withTitleRight(true),
	)

	return m.styles.view.Render(b)
}

func (m Model) headerView() string {

	// s := " Switcher "

	in := Border(
		m.input.View(),
		withWidth(m.input.Width()+1),
		// withExtendSide("right"),
		withTitle("Search"),
		withTitleColor(m.styles.selectedColor),
		withCornorColor(borderColorLight),
		withSideColor(borderColorLight),
		// withCornorChars(BorderCornorRound),
		// withSideChars(BorderSideThin),
	)

	// in := m.input.View()

	// offset := lipgloss.Width(in)

	// line := strings.Repeat(BorderSideThin.Horizontal, max(0, m.viewport.Width()-offset))

	return lipgloss.JoinHorizontal(lipgloss.Center, in)
}

func (m Model) footerView() string {
	s := m.GetSelected()

	title := m.sectionHeader("WORKSPACE DATA")

	wr := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		pathStyle.Render(s.Path),
		s.Branch,
	)

	return wr
}

func (m *Model) listView() string {
	var s strings.Builder

	for _, section := range m.filtered {

		icon := section.Icon + " "

		sh := m.sectionHeader(icon + section.Title)

		s.WriteString(sh + "\n")

		for j, workspace := range section.List {

			itemStr := m.ItemView(workspace)

			// l := "├─ "
			l := ""

			newline := ""
			if len(section.List)-1 == j {
				// l = "└─ "
				l = ""
				newline = "\n"
			}

			s.WriteString(mutedStyle.Render(l) + itemStr + newline)

		}
	}
	return s.String()
}

func (m Model) sectionHeader(title string) string {

	text := sectionTitleStyle.Render("" + title)
	width := m.viewport.Width()

	length := lipgloss.Width(text)
	rem := max(0, width-length)

	sep := mutedStyle.Render(strings.Repeat(BorderSideThin.Horizontal, rem))
	text = text + "" + sep

	return text + ""
}

func (m Model) ItemView(workspace data.Workspace) string {

	str := workspace.Title
	title := lipgloss.JoinHorizontal(lipgloss.Center,
		workspace.Title,
	)

	if m.SelectedId == workspace.Id {
		str = m.styles.listItem.Width(m.width).
			Background(m.styles.selectedColor).
			Render(title)
	} else {
		str = m.styles.listItem.Width(m.width).Render(title)
	}

	return str + "\n"
}

func (m *Model) GetSelected() data.Workspace {
	l := data.MergeSectionWorkspaces(m.filtered)
	if len(l) != 0 {
		filter := data.Find(l, func(w data.Workspace) bool {
			return w.Id == m.SelectedId
		})
		return filter[0]
	}
	return data.Workspace{}
}
