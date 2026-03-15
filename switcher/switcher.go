package switcher

import (
	"image/color"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	input      textinput.Model
	list       []Section
	filtered   []Section
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

	// borderColorLight = lipgloss.Color("1")
	// borderColorDark  = lipgloss.Color("10")

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(borderColorDark).
				PaddingRight(1)

	mutedStyle = lipgloss.NewStyle().Foreground(borderColorDark)

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

	pathStyle = lipgloss.NewStyle().Foreground(borderColorLight)
)

func New(sections []Section, width int, height int) Model {

	i := textinput.New()
	i.Placeholder = "Search…"
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

		widthOffset := 20
		heightOffset := 4
		width := (msg.Width - widthOffset) / 2
		height := (msg.Height - verticalMarginHeight - heightOffset) / 2

		// m.input.SetWidth(width)

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
			newSections := []Section{}
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
	// m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var wrapper string

	if !m.ready {

		wrapper = "\n init..."

	} else {

		wrapper = lipgloss.JoinVertical(lipgloss.Left, m.headerView(), m.viewport.View())

	}

	// t := coloredBorder(wrapper, m.viewport.Width(), m.viewport.Height(), borderColorLight, borderColorDark, "")

	return m.styles.view.Render(wrapper)
}

func (m Model) headerView() string {

	s := " Switcher "

	// in := titleStyle.Render(m.input.View())
	in := coloredBorder(m.input.View(), m.input.Width()+1, 0, borderColorLight, borderColorDark, "right")

	// title := infoStyle.Render("Switcher")

	title := coloredBorder(s, lipgloss.Width(s), 0, borderColorLight, borderColorDark, "left")

	offset := lipgloss.Width(in) + lipgloss.Width(title)

	line := strings.Repeat("─", max(0, m.viewport.Width()-offset))
	// line := strings.Repeat("━", max(0, m.viewport.Width()-offset))

	return lipgloss.JoinHorizontal(lipgloss.Center, in, lineStyle.Render(line), title)
}

func (m Model) footerView() string {
	return ""
}

func (m *Model) listView() string {
	var s strings.Builder

	for _, section := range m.filtered {

		sh := m.sectionHeader(section)

		s.WriteString(sh)

		for j, workspace := range section.List {

			itemStr := m.ItemView(workspace)

			l := "├─ "

			newline := ""
			if len(section.List)-1 == j {
				l = "└─ "
				newline = "\n"
			}

			s.WriteString(mutedStyle.Render(l) + itemStr + newline)

		}
	}
	return s.String()
}

func (m Model) sectionHeader(section Section) string {

	text := sectionTitleStyle.Render("" + section.Title)
	width := m.width

	length := lipgloss.Width(text)
	rem := width - length

	if rem > 0 {
		sep := mutedStyle.Render(strings.Repeat("─", rem))
		text = text + "" + sep
	}

	return text + "\n"
}

func (m Model) ItemView(workspace Workspace) string {
	// home, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// p := strings.ReplaceAll(workspace.Path, home, "~")

	str := workspace.Title
	title := lipgloss.JoinHorizontal(lipgloss.Center,
		workspace.Title,
		// strings.Repeat(" ", 20-lipgloss.Width(workspace.Title)),
		// workspace.Branch,
		// strings.Repeat(" ", 20-lipgloss.Width(workspace.Branch)),
		// pathStyle.Render(p),
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

func coloredBorder(content string, width, height int, cornerColor, sideColor color.Color, side string) string {
	cornerStyle := lipgloss.NewStyle().Foreground(cornerColor)
	sideStyle := lipgloss.NewStyle().Foreground(sideColor)

	// topLeft := cornerStyle.Render("┌")
	// topRight := cornerStyle.Render("┐")
	// bottomLeft := cornerStyle.Render("└")
	// bottomRight := cornerStyle.Render("┘")

	topLeft := cornerStyle.Render("┏")
	topRight := cornerStyle.Render("┓")
	bottomLeft := cornerStyle.Render("┗")
	bottomRight := cornerStyle.Render("┛")

	// horizontal := sideStyle.Render(strings.Repeat("━", width))
	horizontal := sideStyle.Render(strings.Repeat("─", width))

	leftVertical := sideStyle.Render("│")
	// leftVertical := sideStyle.Render("┃")

	if side == "left" {
		leftVertical = sideStyle.Render("┤")
		// leftVertical = sideStyle.Render("┫")
	}

	rightVertical := sideStyle.Render("│")
	// rightVertical := sideStyle.Render("┃")

	if side == "right" {
		rightVertical = sideStyle.Render("├")
		// rightVertical = sideStyle.Render("┣")
	}

	top := topLeft + horizontal + topRight
	bottom := bottomLeft + horizontal + bottomRight

	lines := strings.Split(content, "\n")
	var middle strings.Builder
	for _, line := range lines {
		middle.WriteString(leftVertical + line + rightVertical + "\n")
	}

	return top + "\n" + middle.String() + bottom
}
