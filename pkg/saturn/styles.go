package saturn

import (
	"github.com/charmbracelet/lipgloss"
)

var DefaultStyle = lipgloss.NewStyle()

var styles = map[string]lipgloss.Style{
	"p": DefaultStyle,
}

func Style(tag string) lipgloss.Style {
	if style, ok := styles[tag]; !ok {
		return DefaultStyle
	} else {
		return style
	}
}

func style1(baseStyle lipgloss.Style, style string) lipgloss.Style {
	switch style {
	case "title":
		return baseStyle.Bold(true).Foreground(lipgloss.Color("9"))
	case "highlight":
		return baseStyle.Foreground(lipgloss.Color("5"))
	case "italic", "i":
		return baseStyle.Italic(true)
	case "bold":
		return baseStyle.Bold(true).Foreground(lipgloss.Color("9"))
	case "underline":
		return baseStyle.Underline(true)
	case "p":
		return baseStyle.Foreground(lipgloss.Color("12"))
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return baseStyle.Bold(true).Foreground(lipgloss.Color("9"))
	case "cursor":
		return baseStyle.Reverse(true)
	}
	return baseStyle
}
