package parser

import (
	"strings"

	"github.com/elinx/saturn/pkg/epub"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

func parse(n *html.Node, formater IHtmlFormater) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += parse(c, formater)
	}
	if n.Type == html.ElementNode {
		switch n.Data {
		case "i":
			return formater.I(result)
		case "p":
			return formater.P(result, n.Attr)
		case "title":
			return formater.Title(result)
		case "h1", "h2", "h3", "h4", "h5", "h6":
			return formater.Header(result)
		}
	}
	return result
}

func Parse(content string, formater IHtmlFormater) (string, error) {
	log.Infoln("Enter into parsing of HTML")
	htmlNode, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", err
	}
	return formater.PostProcess(parse(htmlNode, formater)), nil
}

type Segment struct {
	Content string
	Style   string
	Pos     int
}
type Line struct {
	Content  string
	Segments []Segment
	Style    string
}
type Buffer struct {
	Lines []Line
}

type Renderer struct {
	book   *epub.Epub
	buffer *Buffer
}

func New(book *epub.Epub) *Renderer {
	return &Renderer{book: book, buffer: &Buffer{}}
}

func (r *Renderer) Render() error {
	content, err := r.book.GetSpinContent()
	if err != nil {
		return err
	}
	for _, id := range content.Orders {
		htmlContent := content.Contents[id]
		r.Parse(htmlContent)
	}
	return nil
}

func (r *Renderer) Parse(content string) error {
	log.Infoln("Enter into parsing of HTML")
	htmlNode, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return err
	}
	if _, err := r.parse(htmlNode); err != nil {
		return err
	}
	return nil
}

func (r *Renderer) parse(n *html.Node) (*Segment, error) {
	switch n.Type {
	case html.TextNode:
		if len(strings.TrimSpace(n.Data)) == 0 {
			return nil, nil
		}
		return &Segment{n.Data, "", 0}, nil
	case html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if _, err := r.parse(c); err != nil {
				return nil, err
			}
		}
		return nil, nil
	case html.CommentNode:
		return nil, nil
	}
	var segments []Segment
	pos := 0
	contents := []string{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		segment, err := r.parse(c)
		if err != nil {
			return nil, err
		}
		if segment != nil {
			segment.Pos = pos
			segments = append(segments, *segment)
			contents = append(contents, segment.Content)
			pos += len(segment.Content)
		}
	}
	lineContent := strings.Join(contents, "")
	switch n.Data {
	case "head", "html", "body":
		// ignore
	// case "body":
	// 	if len(segments) != 0 {
	// 		r.buffer.Lines = append(r.buffer.Lines, Line{lineContent, segments, ""})
	// 	}
	case "i":
		if n.Parent.Data == "body" {
			r.buffer.Lines = append(r.buffer.Lines, Line{lineContent, segments, n.Data})
			return nil, nil
		}
		return &Segment{lineContent, "i", 0}, nil
	default:
		r.buffer.Lines = append(r.buffer.Lines, Line{lineContent, segments, n.Data})
	}
	return nil, nil
}
