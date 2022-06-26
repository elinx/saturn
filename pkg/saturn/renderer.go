package saturn

import (
	"strings"
	"unicode/utf8"

	"github.com/elinx/saturn/pkg/epub"
	"github.com/elinx/saturn/pkg/util"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
)

type Renderer struct {
	book   *epub.Epub
	buffer *Buffer

	wrapWidth int
}

func NewRender(book *epub.Epub, buffer *Buffer) *Renderer {
	r := &Renderer{
		book:   book,
		buffer: buffer,
	}
	buffer.renderer = r
	return r
}

// Render iterates over the buffer and renders each line to the screen.
func (r *Renderer) Render(width int) {
	lineNumAccum := 0
	r.wrapWidth = width
	for _, line := range r.buffer.Lines {
		content := renderLine1(line)
		lineWraped := util.Wrap(content, width)
		r.buffer.visualLineOffset = append(r.buffer.visualLineOffset, VisualLineIndex(lineNumAccum))
		lineNum := strings.Count(lineWraped, "\n") + 1
		lineNumAccum += lineNum
		visualLines := strings.Split(lineWraped, "\n")
		r.buffer.visualLines = append(r.buffer.visualLines, visualLines...)
	}
}

func (r *Renderer) GetBuffer() *Buffer {
	return r.buffer
}

func (r *Renderer) GetBufferX(line string, vy VisualLineIndex, vx VisualIndex) RuneIndex {
	return RuneIndex(util.LocBeforeWraped(line, r.wrapWidth, int(vx), int(vy)))
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

func renderLine1(line Line) string {
	result := ""
	content := line.Content
	index := ByteIndex(0)
	for len(content) > 0 {
		rune, size := utf8.DecodeRuneInString(content)
		styled := style(string(rune), 0, line.Style)
		for _, s := range line.Segments {
			if s.Pos <= index && s.Pos+ByteIndex(len(s.Content)) > index {
				styled = style(styled, s.Pos, s.Style)
			}
		}
		result += styled
		index += ByteIndex(size)
		content = content[size:]
	}
	return result
}

func (r *Renderer) GetVisualLineNumById(id epub.ManifestId) VisualLineIndex {
	return r.buffer.GetVisualLineNumById(id)
}

func rune2ByteIndex(line string, runeIndex RuneIndex) ByteIndex {
	return ByteIndex(len(string([]rune(line)[:runeIndex])))
}

func (r *Renderer) updateVisualLines(linum BufferLineIndex, line Line) {
	lines := strings.Split(util.Wrap(renderLine1(line), r.wrapWidth), "\n")
	log.Debugf("rendering line %d: %s", linum, lines)
	for i, v := range lines {
		log.Debugf("updating line %d: %s", r.buffer.visualLineOffset[linum]+VisualLineIndex(i), v)
		r.buffer.visualLines[r.buffer.visualLineOffset[linum]+VisualLineIndex(i)] = v
	}
}

func (r *Renderer) MarkPosition(vy VisualLineIndex, vx VisualIndex) {
	y := r.buffer.GetBufferLineNumByVisual(vy)
	x := r.buffer.GetBufferX(y, vy, vx)
	log.Debugf("MarkPosition: %d, %d", y, x)
	line := r.buffer.Lines[y]
	// TODO: show one space if the cursor is at an empty line
	if len(line.Content) > 0 {
		line.Segments = append(line.Segments, Segment{
			Content: string([]rune(line.Content)[x]),
			Style:   "cursor",
			Pos:     rune2ByteIndex(line.Content, x),
		})
		r.buffer.Lines[y] = line
	}
	r.updateVisualLines(y, line)
}

func (r *Renderer) ClearCursorStyles(vy VisualLineIndex) {
	y := r.buffer.GetBufferLineNumByVisual(vy)
	line := r.buffer.Lines[y]
	segments := []Segment{}
	for _, s := range line.Segments {
		if s.Style != "cursor" {
			segments = append(segments, s)
		}
	}
	line.Segments = segments
	r.buffer.Lines[y] = line
	r.updateVisualLines(y, line)
}

func (r *Renderer) MarkInline(vy VisualLineIndex, vxs, vxe VisualIndex) {
	y := r.buffer.GetBufferLineNumByVisual(vy)
	xs := r.buffer.GetBufferX(y, vy, vxs)
	xe := r.buffer.GetBufferX(y, vy, vxe)

	line := r.buffer.Lines[y]
	if len(line.Content) > 0 {
		line.Segments = append(line.Segments, Segment{
			Content: string([]rune(line.Content)[xs : xe+1]),
			Style:   "cursor",
			Pos:     rune2ByteIndex(line.Content, xs),
		})
		r.buffer.Lines[y] = line
	}
	r.updateVisualLines(y, line)
}
