package switcher

import (
	"dashboard/data"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type previewData struct {
	id     string
	data   string
	branch string
	diff   string
}

type Model struct {
	input           textinput.Model
	list            []data.Section
	filtered        []data.Section
	openWorkspaces  []data.Workspace
	previewData     []previewData
	SelectedId      string
	SelectedData    data.Workspace
	width           int
	height          int
	styles          *Styles
	viewport        viewport.Model
	previewViewport viewport.Model
	ready           bool
}

var (
	borderColorLight = lipgloss.Color("#808080")
	borderColorDark  = lipgloss.Color("#505050")

	TitleColor = lipgloss.Color("#ffffff")

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(borderColorDark).
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

func New(sections []data.Section, openWorkspaces []data.Workspace, width int, height int) Model {

	i := textinput.New()
	i.Placeholder = ""
	i.SetVirtualCursor(true)
	i.Focus()
	i.CharLimit = 156
	i.SetWidth(20)
	// i.Prompt = "❯ "
	i.Prompt = ""

	m := Model{
		openWorkspaces: openWorkspaces,
		input:          i,
		list:           sections,
		filtered:       sections,
		width:          width,
		height:         height,
		styles:         DefaultStyles(width+1, height),
	}

	m.setSelected(sections[0].List[0].Id)

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

func (m *Model) setSize(msgWidth, msgHeight int) {

	headerHeight := lipgloss.Height(m.headerView())
	// footerHeight := lipgloss.Height(m.preview())
	verticalMarginHeight := headerHeight

	previewHeaderHeight := lipgloss.Height(m.previewHeader())
	previewFooterHeight := lipgloss.Height(m.previewFooter())

	previewVerticalMarginHeight := previewHeaderHeight + previewFooterHeight

	widthOffset := 2
	heightOffset := 2
	width := (msgWidth - widthOffset) / 2
	height := (msgHeight - verticalMarginHeight - heightOffset)

	previewWidth := (msgWidth - (widthOffset + 3)) / 2
	previewHeight := (msgHeight - previewVerticalMarginHeight - heightOffset)

	m.input.SetWidth(width - 1)

	if !m.ready {

		m.previewViewport = viewport.New(viewport.WithWidth(previewWidth), viewport.WithHeight(previewHeight))
		m.previewViewport.YPosition = previewHeaderHeight

		m.viewport = viewport.New(viewport.WithWidth(width), viewport.WithHeight(height))
		m.viewport.YPosition = headerHeight

		m.viewport.SetContent(m.listView())

		m.previewViewport.SetContent(m.previewContent())

		m.ready = true

	} else {
		m.viewport.SetWidth(width)
		m.viewport.SetHeight(height)

		m.previewViewport.SetWidth(previewWidth)
		m.previewViewport.SetHeight(previewHeight)

		m.previewViewport.SetContent(m.previewContent())

		m.viewport.SetContent(m.listView())

		m.width = width
		m.height = height
	}

}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	switch msg := msg.(type) {

	case data.UpdateStateMsg:
		m.setSize(msg.Width, msg.Height)

	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		switch msg.String() {

		case "up":
			newId := m.moveUp()
			if newId != "" {
				m.setSelected(newId)
			}

		case "down":
			newId := m.moveDown()
			if newId != "" {
				m.setSelected(newId)

			}

		case "enter":
			l := data.MergeSectionWorkspaces(m.filtered)
			if len(l) != 0 {
				filter := data.Find(l, func(w data.Workspace) bool {
					return w.Id == m.SelectedId
				})
				selected := filter[0].Path
				enc := setUserVar(selected)
				fmt.Print(enc)
				time.Sleep(500 * time.Millisecond)
				return m, tea.Quit
			}

		case "esc":
			return m, tea.Quit

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
				id := m.filtered[0].List[0].Id
				m.setSelected(id)
			}

			m.viewport.SetContent(m.listView())

			return m, cmd
		}
	}

	var cmd tea.Cmd
	return m, cmd
}

func setUserVar(value string) string {
	key := "go-cli"

	encoded := base64.StdEncoding.EncodeToString([]byte(value))

	escape := fmt.Sprintf("\033]1337;SetUserVar=%s=%s\007", key, encoded)

	return escape
}

func (m *Model) setSelected(newId string) {
	m.SelectedId = newId
	m.scrollToSelected()

	item := m.GetSelected(newId)

	m.SelectedData = item

	// f := data.Find(m.previewData, func(p previewData) bool {
	// 	return p.id == newId
	// })

	// if len(f) == 0 {
	//
	// 	branch := data.GetCmdOut(item.Path, "git", "branch", "--show-current")
	// 	diff := data.Diff(item.Path)
	// 	data := listDirs(item.Path)
	//
	// 	newData := previewData{id: newId, data: data, branch: branch, diff: diff}
	// 	m.previewData = append(m.previewData, newData)
	// }

	m.viewport.SetContent(m.listView())
	// m.previewViewport.SetContent(m.previewContent())
}

func (m Model) View() string {
	var wrapper string

	if !m.ready {

		wrapper = "\n init..."

	} else {

		WorkspaceBorderStyle := Border(
			WithWidth(m.viewport.Width()),
			WithHeight(m.viewport.Height()),
			WithTitleColor(TitleColor),
			WithTitle("Workspaces"),
		)

		previewBorderStyle := Border(
			WithWidth(m.previewViewport.Width()),
			WithHeight(m.previewViewport.Height()),
			WithTitleColor(TitleColor),
			WithTitle("Preview"),
		)

		left := lipgloss.JoinVertical(lipgloss.Left, m.headerView(), WorkspaceBorderStyle.Render(m.viewport.View()))

		wrapper = lipgloss.JoinHorizontal(lipgloss.Left, left, previewBorderStyle.Render(m.previewView()))

	}

	return m.styles.view.Render(wrapper)
}

func (m Model) headerView() string {

	InputBorderStyle := Border(
		WithWidth(m.input.Width()+1),
		WithTitle("Search"),
		WithTitleColor(TitleColor),
	)

	return lipgloss.JoinHorizontal(lipgloss.Center, InputBorderStyle.Render(m.input.View()))
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

	cursor := "  "

	title := lipgloss.JoinHorizontal(lipgloss.Left,
		workspace.Title,
	)

	if m.SelectedId == workspace.Id {
		cursor = m.styles.listItem.Foreground(m.styles.selectedColor).Render(" ")
		str = m.styles.listItem.Width(m.width).
			// Background(m.styles.selectedColor).
			Render(cursor + title)
	} else {
		str = m.styles.listItem.Width(m.width).Render(cursor + title)
	}

	return str + "\n"
}

func (m *Model) GetSelected(id string) data.Workspace {
	l := data.MergeSectionWorkspaces(m.filtered)
	if len(l) != 0 {
		filter := data.Find(l, func(w data.Workspace) bool {
			return w.Id == id
		})
		return filter[0]
	}
	return data.Workspace{}
}
