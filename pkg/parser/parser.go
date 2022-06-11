package parser

import (
	"strings"

	"github.com/elinx/saturn/pkg/epub"
	"github.com/muesli/reflow/wrap"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

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

	// The position of each block of the spine in the Lines
	BlockPos map[epub.ManifestId]int
}

type Renderer struct {
	book    *epub.Epub
	buffer  *Buffer
	offsets []int
}

func New(book *epub.Epub) *Renderer {
	return &Renderer{book: book, buffer: &Buffer{
		Lines:    []Line{},
		BlockPos: make(map[epub.ManifestId]int),
	}}
}

func (r *Renderer) Render(width int) string {
	var lines []string
	lineNumAccum := 0
	for _, line := range r.buffer.Lines {
		lineWraped := wrap.String(line.Content, width)
		r.offsets = append(r.offsets, lineNumAccum)
		lineNum := strings.Count(lineWraped, "\n") + 1
		lineNumAccum += lineNum
		lines = append(lines, lineWraped)
	}
	return strings.Join(lines, "\n")
}

// renderWrap wraps the content of the line with the given width
func (r *Renderer) renderWrap(line string, width int) (string, int) {
	lineWraped := wrap.String(line, width)
	lineNum := strings.Count(lineWraped, "\n") + 1
	return lineWraped, lineNum
}

func (r *Renderer) GetPos(id epub.ManifestId) int {
	return r.buffer.BlockPos[id]
}

func (r *Renderer) GetVisualPos(id epub.ManifestId) int {
	return r.offsets[r.buffer.BlockPos[id]]
}

// Parse iterates over the spine and parses each HTML file
func (r *Renderer) Parse() error {
	content, err := r.book.GetSpinContent()
	if err != nil {
		return err
	}
	for _, id := range content.Orders {
		htmlContent := content.Contents[id]
		r.buffer.BlockPos[id] = len(r.buffer.Lines)
		r.parse1(htmlContent)
	}
	return nil
}

func (r *Renderer) parse1(content string) error {
	log.Infoln("Enter into parsing of HTML")
	htmlNode, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return err
	}
	if _, err := r.parse2(htmlNode); err != nil {
		return err
	}
	return nil
}

func (r *Renderer) parse2(n *html.Node) (*Segment, error) {
	switch n.Type {
	case html.TextNode:
		if len(strings.TrimSpace(n.Data)) == 0 {
			return nil, nil
		}
		return &Segment{n.Data, "", 0}, nil
	case html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if _, err := r.parse2(c); err != nil {
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
		segment, err := r.parse2(c)
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
	case "head", "html", "body", "link":
		// ignore
	case "svg", "image", "img":
		// TODO: support image display
	case "style":
		// TODO: support inline style
	case "i", "b", "strong", "span", "em":
		if n.Parent.Data == "body" {
			r.buffer.Lines = append(r.buffer.Lines, Line{lineContent, segments, n.Data})
			return nil, nil
		}
		return &Segment{lineContent, n.Data, 0}, nil
	default:
		r.buffer.Lines = append(r.buffer.Lines, Line{lineContent, segments, n.Data})
	}
	return nil, nil
}
