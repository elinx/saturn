package html_parser

import (
	"strings"

	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type C struct {
	c     rune
	X, Y  int
	style termenv.Style
}

type Line []C

func (line Line) TrimLeadingSpace() Line {
	var j int
	for i := range line {
		if line[i].c == ' ' || line[i].c == '\t' {
			j++
		} else {
			break
		}
	}
	if j >= len(line) {
		j = len(line) - 1
	}
	return line[j:]
}

func (line Line) TrimTrailingSpace() Line {
	// The last character is always a new line
	if len(line) <= 1 {
		return line
	}
	var j int = len(line) - 2
	for i := len(line) - 2; i >= 0; i-- {
		if line[i].c == ' ' || line[i].c == '\t' {
			line[i].c = '\n'
			j--
		} else {
			break
		}
	}
	if j < 0 {
		j = 0
	}
	return line[:j+2]
}

func (line Line) Render() string {
	var builder strings.Builder
	line = line.TrimLeadingSpace()
	line = line.TrimTrailingSpace()
	for _, c := range line {
		builder.WriteString(string(c.style.Styled(string(c.c))))
	}
	return builder.String()
}

func string2Cs(str string, style termenv.Style) Content {
	var cs Line
	for _, c := range str {
		cs = append(cs, C{c, 0, 0, style})
	}
	return Content{cs}
}

type Content []Line

func (c Content) Render() string {
	// TODO: skip continuaus white spaces
	var builder strings.Builder
	for _, line := range c {
		builder.WriteString(line.Render())
	}
	return builder.String()
}

type RuneFormater interface {
	GetDefaultStyle() termenv.Style
	I(Content) Content
	P(Content) Content
	Title(Content) Content
	Header(Content) Content
}

func (f formater) GetDefaultStyle() termenv.Style {
	return termenv.String()
}

func (f formater) I(c Content) Content {
	for i := range c {
		for j := range c[i] {
			c[i][j].style = c[i][j].style.Italic()
		}
	}
	return c
}

func (f formater) P(c Content) Content {
	if len(c) == 0 {
		return c
	}
	// Join lines
	firstLine := c[0]
	for line := range c[1:] {
		firstLine = append(firstLine, c[line+1]...)
	}
	// Add new line to the end of the content only if it is not already there.
	if firstLine[len(firstLine)-1].c != '\n' {
		firstLine = append(firstLine, C{'\n', 0, 0, f.GetDefaultStyle()})
	}
	return Content{firstLine}
}

func (f formater) Title(c Content) Content {
	for i := range c {
		for j := range c[i] {
			c[i][j].style.Bold()
		}
	}
	return c
}

func (f formater) Header(c Content) Content {
	for i := range c {
		for j := range c[i] {
			c[i][j].style.Bold()
		}
	}
	return c
}

type formater struct{}

func parse(node *html.Node, formater RuneFormater) Content {
	if node.Type == html.TextNode {
		text := strings.TrimSpace(node.Data)
		if text == "" {
			return nil
		}
		// TODO: trim extra trailing or leading spaces in text
		return string2Cs(node.Data, formater.GetDefaultStyle())
	}
	var result Content
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, parse(c, formater)...)
	}
	if node.Type == html.ElementNode {
		switch node.Data {
		case "i":
			result = formater.I(result)
		case "p":
			result = formater.P(result)
		case "title":
			result = formater.Title(result)
		case "h1", "h2", "h3", "h4", "h5", "h6":
			result = formater.Header(result)
		}
	}
	return result
}

func Parse(content string, formater RuneFormater) (Content, error) {
	log.Infoln("Enter into parsing of HTML")
	htmlNode, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	return parse(htmlNode, formater), nil
}
