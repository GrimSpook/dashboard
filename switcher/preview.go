package switcher

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	// "github.com/epilande/go-devicons"
)

var (
	preStyle = lipgloss.NewStyle().PaddingLeft(0)

	titleS = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)
)

func (m Model) previewView() string {

	wr := lipgloss.JoinVertical(
		lipgloss.Left,
		m.previewHeader(),
		preStyle.Render(m.previewViewport.View()),
		m.previewFooter(),
	)

	return wr
}

func (m Model) previewHeader() string {

	var str strings.Builder

	for _, data := range m.previewData {
		if m.SelectedId == data.id {

			b := "" + titleS.Render(strings.Trim(data.branch, "\n")) + ""

			if len(data.branch) == 0 {
				b = ""
			}

			s := lipgloss.JoinVertical(lipgloss.Left,
				"",
				b,
				"",
			)

			str.WriteString(s)

		}
	}

	return str.String()
}

func (m Model) previewContent() string {
	for _, data := range m.previewData {
		if m.SelectedId == data.id {

			if len(data.branch) == 0 {
				return data.data
			}

			di := data.diff
			if len(data.diff) == 0 && len(data.branch) != 0 {
				di = "No changes\n"
			}

			return lipgloss.JoinVertical(
				lipgloss.Left,
				di,
				data.data,
			)
		}
	}
	return ""
}

func (m Model) previewFooter() string {

	p := pathStyle.Render(m.SelectedData.Path)

	str := lipgloss.JoinVertical(lipgloss.Left,
		"",
		// mutedStyle.Render(strings.Repeat(BorderSideThin.Horizontal, m.previewViewport.Width())),
		p,
	)

	return str
}

func listDirs(path string) string {

	nameStyle := lipgloss.NewStyle()
	// iconStyle := lipgloss.NewStyle()

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}

	var str []string

	base := nameStyle.Foreground(lipgloss.Color("4")).Bold(true)

	str = append(str, base.Render(filepath.Base(path)))

	for i, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Fatalln(err)
		}

		// fileStyle := devicons.IconForInfo(info)

		color := lipgloss.Color("15")

		if entry.Name()[:1] == "." {
			color = lipgloss.Color("8")
		} else if entry.IsDir() {
			color = lipgloss.Color("4")
		}

		// fileColor := lipgloss.Color("1")

		name := nameStyle.Foreground(color).Bold(info.IsDir()).Render(entry.Name())
		// icon := iconStyle.Foreground(color).Render(fileStyle.Icon)

		l := "├─ "

		if len(entries)-1 == i {
			l = "└─ "
		}

		s := lipgloss.JoinHorizontal(
			lipgloss.Left,
			mutedStyle.Render(l),
			// icon,
			"",
			name,
		)

		if entry.Name()[:1] != "." {

			str = append(str, s)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, str...)

}
