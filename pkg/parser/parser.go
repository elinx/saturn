package parser

import (
	"strings"
	"unicode/utf8"

	"github.com/elinx/saturn/pkg/epub"
	"github.com/elinx/saturn/pkg/util"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

// RuneIndex returns the index of the rune in the given string
type RuneIndex int

// ByteIndex returns the index of the byte in the given string
type ByteIndex int

// VisualIndex returns the visual index of the screen position
type VisualIndex int

type Segment struct {
	Content string
	Style   string
	Pos     ByteIndex
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
	book   *epub.Epub
	buffer *Buffer

	wrapWidth int

	// lineYOffsets is the offset of each line in the buffer after
	// being rendered to the screen. It is used to calculate the
	// position of each rune in the line.
	lineYOffsets []int
}

func New(book *epub.Epub) *Renderer {
	return &Renderer{book: book, buffer: &Buffer{
		Lines:    []Line{},
		BlockPos: make(map[epub.ManifestId]int),
	}}
}

// Render iterates over the buffer and renders each line to the screen.
func (r *Renderer) Render(width int) string {
	var lines []string
	lineNumAccum := 0
	r.wrapWidth = width
	for _, line := range r.buffer.Lines {
		content := renderLine1(line)
		lineWraped := util.Wrap(content, width)
		r.lineYOffsets = append(r.lineYOffsets, lineNumAccum)
		lineNum := strings.Count(lineWraped, "\n") + 1
		lineNumAccum += lineNum
		lines = append(lines, lineWraped)
	}
	return strings.Join(lines, "\n")
}

// style returns the style of the given content with the given style
// TODO: use a style sheet
func style(content string, pos ByteIndex, style string) string {
	switch style {
	case "title":
		content = termenv.String(content).Foreground(termenv.ANSIBrightRed).Bold().String()
	case "highlight":
		content = termenv.String(content).Underline().String()
	case "italic", "i":
		content = termenv.String(content).Italic().String()
	case "bold":
		content = termenv.String(content).Bold().String()
	case "underline":
		content = termenv.String(content).Underline().String()
	case "p":
		content = termenv.String(content).Foreground(termenv.ANSICyan).String()
	case "h1", "h2", "h3", "h4", "h5", "h6":
		content = termenv.String(content).Foreground(termenv.ANSIBrightRed).Bold().String()
	case "cursor":
		content = termenv.String(content).Reverse().Blink().String()
	default:
		content = termenv.String(content).String()
	}
	return content
}

func renderLine(line Line) string {
	result := ""
	cursors := []*Segment{}
	for _, s := range line.Segments {
		result += style(s.Content, s.Pos, s.Style)
		if s.Style == "cursor" {
			cursors = append(cursors, &s)
		}
	}
	for _, s := range cursors {
		width := len(s.Content)
		result = line.Content[:s.Pos] +
			termenv.String(string(line.Content[s.Pos:s.Pos+ByteIndex(width)])).Reverse().Blink().String() +
			line.Content[s.Pos+ByteIndex(width):]
	}
	return style(result, 0, line.Style)
}

func renderLine1(line Line) string {
	result := ""
	content := line.Content
	index := ByteIndex(0)
	for len(content) > 0 {
		rune, size := utf8.DecodeRuneInString(content)
		styled := string(rune)
		for _, s := range line.Segments {
			if s.Pos <= index && s.Pos+ByteIndex(len(s.Content)) > index {
				styled = style(styled, s.Pos, s.Style)
			}
		}
		result += styled
		index += ByteIndex(size)
		content = content[size:]
	}
	return style(result, 0, line.Style)
}

// renderWrap wraps the content of the line with the given width
func (r *Renderer) renderWrap(line string, width int) (string, int) {
	lineWraped := util.Wrap(line, width)
	lineNum := strings.Count(lineWraped, "\n") + 1
	return lineWraped, lineNum
}

func (r *Renderer) GetPos(id epub.ManifestId) int {
	return r.buffer.BlockPos[id]
}

func (r *Renderer) GetVisualYPos(id epub.ManifestId) int {
	return r.lineYOffsets[r.buffer.BlockPos[id]]
}

func (r *Renderer) GetVisualYPos1(line int) int {
	return r.lineYOffsets[line]
}

func rune2ByteIndex(line string, runeIndex RuneIndex) ByteIndex {
	return ByteIndex(len(string([]rune(line)[:runeIndex])))
}

func (r *Renderer) MarkPosition(lineNum int, x RuneIndex) {
	// lineNum = 0
	// x = 0
	line := r.buffer.Lines[lineNum]
	line.Segments = append(line.Segments, Segment{
		Content: string([]rune(line.Content)[x]),
		Style:   "cursor",
		Pos:     rune2ByteIndex(line.Content, x),
	})
	r.buffer.Lines[lineNum] = line
}

func (r *Renderer) GetOriginYPos(visualLineNum int) int {
	for i, v := range r.lineYOffsets {
		if v == visualLineNum {
			return i
		} else if v > visualLineNum {
			return i - 1
		}
	}
	return len(r.lineYOffsets) - 1
}

func (r *Renderer) GetOriginXPos(originLineNum int, visualXPos, visualYPos int) RuneIndex {
	line := r.buffer.Lines[originLineNum].Content
	return RuneIndex(util.LocBeforeWraped(line, r.wrapWidth, visualXPos, visualYPos))
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
	pos := ByteIndex(0)
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
			pos += ByteIndex(len(segment.Content))
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
