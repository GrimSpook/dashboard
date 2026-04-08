package tracker

import (
	"dashboard/data"
	sw "dashboard/switcher"
	"log"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	mutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	resStyle      = lipgloss.NewStyle()
	ellipsis      = lipgloss.NewStyle().Bold(true).Render("… ")
)

// type item struct {
// 	company   string
// 	email     string
// 	responded bool
// 	url       string
// 	updatedAt string
// }

type editItem struct {
	value string
	key   string
}

type Model struct {
	list       []data.Company
	editList   []editItem
	viewport   viewport.Model
	cursor     int
	editCursor int
	editInput  bool
	width      int
	height     int
	ready      bool
	edit       bool
	input      textinput.Model
}

func New() Model {
	i := textinput.New()
	i.Placeholder = ""
	i.SetVirtualCursor(true)
	i.Focus()
	i.CharLimit = 156
	i.SetWidth(20)
	i.Prompt = ""

	l, err := data.LoadCompanies()
	if err != nil {
		log.Fatalln(err)
	}

	return Model{
		list:   l,
		cursor: 0,
		input:  i,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) setSize(msgWidth, msgHeight int) {

	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	widthOffset := 2
	heightOffset := 2
	width := (msgWidth - widthOffset)
	height := (msgHeight - verticalMarginHeight - heightOffset)

	if m.ready {

		m.viewport = viewport.New(viewport.WithWidth(width), viewport.WithHeight(height))
		m.viewport.YPosition = headerHeight

		m.viewport.SetContent(m.render())

		m.ready = true

	} else {

		m.viewport.SetWidth(width)
		m.viewport.SetHeight(height)

		m.viewport.SetContent(m.render())

	}

}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	switch msg := msg.(type) {

	case data.UpdateStateMsg:
		m.setSize(msg.Width, msg.Height)

	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		if !m.edit {

			switch msg.String() {

			case "ctrl+n":
				m.list = append(m.list,
					data.Company{
						Name:      "<name>",
						Email:     "<email>",
						Responded: false,
						Url:       "<url>",
						UpdatedAt: time.Now().Format("2006-01-02 15:04"),
					},
				)
				m.cursor = len(m.list) - 1
				m.viewport.SetContent(m.render())

			case "j":
				if m.cursor < len(m.list)-1 {
					m.cursor++
					m.scrollToSelected()
				}
				m.viewport.SetContent(m.render())

			case "k":
				if m.cursor > 0 {
					m.cursor--

					m.scrollToSelected()
				}
				m.viewport.SetContent(m.render())

			case "e":
				m.setEditList()
				m.edit = true
				m.viewport.SetContent(m.render())

			case "esc":
				return m, tea.Quit
			}

		} else {

			if !m.editInput {

				switch msg.String() {

				case "j":
					if m.editCursor < len(m.editList)-1 {
						m.editCursor++
					}

				case "k":
					if m.editCursor > 0 {
						m.editCursor--
					}

				case "e":
					if m.editList[m.editCursor].key != "save" {
						m.input.SetWidth(35)
						m.input.Placeholder = m.editList[m.editCursor].value
						m.input.SetValue(m.editList[m.editCursor].value)
						m.editInput = true
					}

				case "enter":
					if m.editList[m.editCursor].key == "save" {
						m.edit = false
						data.SaveCompanies(m.list)
						m.editCursor = 0
					}

				case "esc":
					m.edit = false
					m.viewport.SetContent(m.render())

				}

			} else {

				switch msg.String() {

				case "esc":
					m.input.SetValue("")
					m.editInput = false

				case "enter":
					m.updateItem()
					m.input.SetValue("")
					m.viewport.SetContent(m.render())
					// m.edit = false
					m.editInput = false
					// data.SaveCompanies(m.list)

				default:

					var cmd tea.Cmd
					m.input, cmd = m.input.Update(msg)

					m.viewport.SetContent(m.render())

					return m, cmd

				}

			}

		}

	}

	var cmd tea.Cmd
	return m, cmd
}

func (m *Model) scrollToSelected() {
	lineNumber := m.cursor

	viewportHeight := m.viewport.Height()
	currentOffset := m.viewport.YOffset()

	scrollPadding := 4

	itemTop := lineNumber - scrollPadding
	itemBottom := lineNumber + scrollPadding

	if itemTop < currentOffset {
		m.viewport.SetYOffset(itemTop)
	}

	if itemBottom > currentOffset+viewportHeight {
		m.viewport.SetYOffset(itemBottom - viewportHeight)
	}
}

func (m Model) View() string {

	prv := m.viewport.View()

	viewBorderStyle := sw.Border(
		sw.WithWidth(m.viewport.Width()),
		sw.WithHeight(m.viewport.Height()),
	)

	if m.edit {
		editViewRender := m.editView()

		w := lipgloss.Width(editViewRender)
		h := lipgloss.Height(editViewRender)

		x := (m.viewport.Width() / 2) - w/2
		y := (m.viewport.Height() / 2) - h/2

		prv = PlaceOverlay(x, y, editViewRender, m.viewport.View())
	}

	s := lipgloss.JoinVertical(
		lipgloss.Left,
		m.headerView(),
		prv,
		m.footerView(),
	)

	return viewBorderStyle.Render(s)
}

func (m *Model) setEditList() {

	m.editList = []editItem{
		{value: m.list[m.cursor].Name, key: "name"},
		{value: m.list[m.cursor].Email, key: "email"},
		{value: m.list[m.cursor].Url, key: "url"},
		{value: "Save", key: "save"},
	}

}

func (m *Model) updateItem() {

	itemMap := make(map[string]editItem)
	for _, im := range m.editList {
		itemMap[im.key] = im
	}

	key := m.editList[m.editCursor].key

	switch key {
	case "name":

		m.list[m.cursor].Name = m.input.Value()
		m.editList[m.editCursor].value = m.input.Value()

	case "email":

		m.list[m.cursor].Email = m.input.Value()

		m.editList[m.editCursor].value = m.input.Value()

	case "url":

		m.list[m.cursor].Url = "https://www." + m.input.Value()

		m.editList[m.editCursor].value = "https://www." + m.input.Value()

	}

	m.list[m.cursor].UpdatedAt = time.Now().Format("2006-01-02 15:04")

}

var (
	btnStyle = lipgloss.NewStyle().Padding(0, 1)
)

func (m Model) editView() string {

	var renderStr []string

	cat := []string{
		"Name",
		"Email",
		"Url",
		"save",
	}

	for i, str := range m.editList {

		s := str.value
		color := lipgloss.Color("240")
		cornor := sw.BorderCornorThin

		if m.editCursor == i {
			s = str.value

			if m.editInput {
				s = m.input.View()
			}

			if str.key == "save" {
				s = str.value
			}

			color = lipgloss.Color("5")
			cornor = sw.BorderCornorThick
		}

		inputBorderStyle := sw.Border(
			sw.WithWidth(36),
			sw.WithTitle(cat[i]),
			sw.WithCornorChars(cornor),
			sw.WithCornorColor(color),
			sw.WithSideColor(color),
			sw.WithTitleBackgroundColor(color),
			sw.WithTitleSide("Left"),
		)

		if str.key == "save" {
			btnRender := btnStyle.Render(s)

			if m.editCursor == i {

				btnRender = btnStyle.Background(lipgloss.Color("5")).Render(s)

			}

			renderStr = append(renderStr, btnRender)
		} else {
			renderStr = append(renderStr, inputBorderStyle.Render(s))
		}

	}

	strCombind := lipgloss.JoinVertical(
		lipgloss.Left,
		renderStr...,
	)

	finalString := lipgloss.NewStyle().Width(40).Padding(1, 1).Render(strCombind)

	editViewBorderStyle := sw.Border(
		sw.WithWidth(40),
		sw.WithHeight(30),
		sw.WithTitle("Edit"),
	)

	return editViewBorderStyle.Render(finalString)
}

func (m Model) formatCatagory(name string, isLast bool) string {
	viewSplit := m.viewport.Width() / 4

	text := name
	if len(name) > viewSplit {
		v := max(0, viewSplit-2)
		text = name[:v] + ellipsis
	}

	s := lipgloss.NewStyle().Width(viewSplit).Render(text)
	if isLast {
		vi := m.viewport.Width() - viewSplit*4
		s = lipgloss.NewStyle().Width(viewSplit + vi).Render(text)
	}

	return s
}

func (m Model) headerView() string {
	cat := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.formatCatagory("Name", false),
		m.formatCatagory("Email", false),
		m.formatCatagory("Updated at", false),
		m.formatCatagory("Sent Email", true),
		// m.formatCatagory("Url"),
	)
	s := lipgloss.JoinVertical(
		lipgloss.Left,
		mutedStyle.Render(cat),
		mutedStyle.Render(strings.Repeat(sw.BorderSideThin.Horizontal, max(0, m.viewport.Width()))),
	)
	return s
}

func (m Model) rowView(item data.Company, index int) string {

	res := "No"
	if item.Responded {
		res = "Yes"
	}

	s := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.formatCatagory(item.Name, false),
		m.formatCatagory(item.Email, false),
		m.formatCatagory(item.UpdatedAt, false),
		m.formatCatagory(res, true),
	)

	if m.cursor == index {
		return selectedStyle.Render(s)
	}

	return s
}

func (m Model) render() string {

	var str []string

	for i, item := range m.list {

		s := m.rowView(item, i)

		str = append(str, s)

	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		str...,
	)
}

func (m Model) footerView() string {

	s := lipgloss.JoinVertical(
		lipgloss.Left,
		mutedStyle.Render(strings.Repeat(sw.BorderSideThin.Horizontal, max(0, m.viewport.Width()))),
		mutedStyle.Render(m.list[m.cursor].Url),
	)

	return s
}
