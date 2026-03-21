package switcher

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	BorderSideThin = BorderSideChars{
		Horizontal:     "─",
		Vertical:       "│",
		ConnectorRight: "├",
		ConnectorLeft:  "┤",
	}

	BorderSideThick = BorderSideChars{
		Horizontal:     "━",
		Vertical:       "┃",
		ConnectorRight: "┣",
		ConnectorLeft:  "┫",
	}

	BorderSideDouble = BorderSideChars{
		Horizontal:     "═",
		Vertical:       "║",
		ConnectorRight: "╠",
		ConnectorLeft:  "╣",
	}

	BorderCornorRound = BorderCornorChars{
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	BorderCornorDouble = BorderCornorChars{
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
	}

	BorderCornorThin = BorderCornorChars{
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	BorderCornorThick = BorderCornorChars{
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}
)

type BorderCornorChars struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
}

type BorderSideChars struct {
	Horizontal     string
	Vertical       string
	ConnectorRight string
	ConnectorLeft  string
}

type BorderStyle struct {
	CornerColor   color.Color
	SideColor     color.Color
	TitleColor    color.Color
	ConnectedSide string
	Title         string
	TitleRight    bool
	SideChar      BorderSideChars
	CornorChar    BorderCornorChars
	Width         int
	Height        int
	TitleBold     bool
}

type BorderOption func(*BorderStyle)

func withCornorColor(c color.Color) BorderOption {
	return func(bs *BorderStyle) { bs.CornerColor = c }
}

func withSideColor(c color.Color) BorderOption {
	return func(bs *BorderStyle) { bs.SideColor = c }
}

func withTitleColor(c color.Color) BorderOption {
	return func(bs *BorderStyle) { bs.TitleColor = c }
}

func withTitle(title string) BorderOption {
	return func(bs *BorderStyle) { bs.Title = title }
}

func withTitleRight(right bool) BorderOption {
	return func(bs *BorderStyle) { bs.TitleRight = right }
}

func withExtendSide(side string) BorderOption {
	return func(bs *BorderStyle) { bs.ConnectedSide = side }
}

func withWidth(w int) BorderOption {
	return func(bs *BorderStyle) { bs.Width = w }
}

func withHeight(h int) BorderOption {
	return func(bs *BorderStyle) { bs.Height = h }
}

func withCornorChars(char BorderCornorChars) BorderOption {
	return func(bs *BorderStyle) { bs.CornorChar = char }
}

func withSideChars(char BorderSideChars) BorderOption {
	return func(bs *BorderStyle) { bs.SideChar = char }
}

func withTitleBold(bold bool) BorderOption {
	return func(bs *BorderStyle) { bs.TitleBold = bold }
}

func Border(content string, opts ...BorderOption) string {

	style := &BorderStyle{
		CornerColor:   lipgloss.Color("#ffffff"),
		SideColor:     lipgloss.Color("#ffffff"),
		TitleColor:    lipgloss.Color("#ffffff"),
		SideChar:      BorderSideThin,
		CornorChar:    BorderCornorThick,
		TitleRight:    false,
		Width:         80,
		Height:        0,
		ConnectedSide: "",
		Title:         "",
		TitleBold:     true,
	}

	for _, opt := range opts {
		opt(style)
	}

	return renderBorder(content, style)

}

func renderBorder(content string, opt *BorderStyle) string {
	cornerStyle := lipgloss.NewStyle().Foreground(opt.CornerColor)
	sideStyle := lipgloss.NewStyle().Foreground(opt.SideColor)
	titleStyle := lipgloss.NewStyle().Foreground(opt.TitleColor)

	topBorder := buildTopBorder(opt, cornerStyle, sideStyle, titleStyle)

	bottomBorder := cornerStyle.Render(opt.CornorChar.BottomLeft) +
		sideStyle.Render(strings.Repeat(opt.SideChar.Horizontal, opt.Width)) +
		cornerStyle.Render(opt.CornorChar.BottomRight)

	middle := buildMiddleContent(content, opt, sideStyle)

	var result strings.Builder
	result.WriteString(topBorder)
	result.WriteString("\n")
	result.WriteString(middle)
	result.WriteString(bottomBorder)

	return result.String()
}

func buildTopBorder(opt *BorderStyle, cornorStyle, sideStyle, titleStyle lipgloss.Style) string {
	topLeft := cornorStyle.Render(opt.CornorChar.TopLeft)
	topRight := cornorStyle.Render(opt.CornorChar.TopRight)

	if opt.Title == "" {
		horizontal := sideStyle.Render(strings.Repeat(opt.SideChar.Horizontal, opt.Width))
		return topLeft + horizontal + topRight
	}

	titleLength := lipgloss.Width(opt.Title)
	availableSpace := opt.Width - titleLength - 4
	// leftDashes := 1
	rightDashes := availableSpace - 1

	if opt.TitleRight {

		horizontal := sideStyle.Render(strings.Repeat(opt.SideChar.Horizontal, max(1, rightDashes))+
			opt.SideChar.ConnectorLeft) +
			" " +
			titleStyle.Render(opt.Title) +
			" " +
			sideStyle.Render(opt.SideChar.ConnectorRight+opt.SideChar.Horizontal)

		return topLeft + horizontal + topRight
	}

	horizontal := sideStyle.Render(opt.SideChar.Horizontal+opt.SideChar.ConnectorLeft) +
		" " +
		titleStyle.Render(opt.Title) +
		" " +
		sideStyle.Render(opt.SideChar.ConnectorRight+strings.Repeat(opt.SideChar.Horizontal, max(1, rightDashes)))

	// horizontal := sideStyle.Render(strings.Repeat(opt.SideChar.Horizontal, max(1, leftDashes))) +
	// 	" " +
	// 	titleStyle.Render(opt.Title) +
	// 	" " +
	// 	sideStyle.Render(strings.Repeat(opt.SideChar.Horizontal, max(1, rightDashes)))

	return topLeft + horizontal + topRight
}

func buildMiddleContent(content string, opt *BorderStyle, sideStyle lipgloss.Style) string {
	leftVertical := sideStyle.Render(opt.SideChar.Vertical)
	rightVertical := sideStyle.Render(opt.SideChar.Vertical)

	// Apply side modifiers
	if opt.ConnectedSide == "left" {
		leftVertical = sideStyle.Render(opt.SideChar.ConnectorLeft)
	}
	if opt.ConnectedSide == "right" {
		rightVertical = sideStyle.Render(opt.SideChar.ConnectorRight)
	}

	lines := strings.Split(content, "\n")
	var result strings.Builder

	for _, line := range lines {

		result.WriteString(leftVertical)
		result.WriteString(line)
		result.WriteString(rightVertical)
		result.WriteString("\n")

	}

	return result.String()

}
