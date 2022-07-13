package saturn

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/elinx/saturn/pkg/epub"
	"github.com/elinx/saturn/pkg/util"
	log "github.com/sirupsen/logrus"
)

type Renderer struct {
	book   *epub.Epub
	buffer *Buffer

	wrapWidth  int
	linumWidth int
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
	r.linumWidth = len(strconv.Itoa(len(r.buffer.Lines)))
	r.wrapWidth = width - r.linumWidth
	for linum := range r.buffer.Lines {
		r.buffer.visualLineOffset = append(r.buffer.visualLineOffset, VisualLineIndex(lineNumAccum))
		visualLines := r.RenderLine(BufferLineIndex(linum))
		lineNumAccum += len(visualLines)
		r.buffer.visualLines = append(r.buffer.visualLines, visualLines...)
	}
}

func (r *Renderer) RenderLinum(linum BufferLineIndex) string {
	line := strconv.Itoa(int(linum))
	line = strings.Repeat(" ", r.linumWidth-len(line)) + line
	return linumStyle.SetString(line).String()
}

func (r *Renderer) RenderEmptyLinum() string {
	return linumStyle.SetString(strings.Repeat(" ", r.linumWidth)).String()
}

func (r *Renderer) RenderLine(linum BufferLineIndex) []VisualLine {
	emptyLinum := r.RenderEmptyLinum()
	rlinum := r.RenderLinum(linum)
	renderedLine := ""
	line := r.buffer.Lines[linum]
	content := line.Content
	index := ByteIndex(0)
	runes := []VisualRune{}
	for len(content) > 0 {
		rune, size := utf8.DecodeRuneInString(content)
		styled := DefaultStyle.SetString(string(rune))
		styled = style1(styled, line.Style)
		for _, s := range line.Segments {
			if s.Pos <= index && s.Pos+ByteIndex(len(s.Content)) > index {
				styled = style1(styled, s.Style)
			}
		}
		renderedLine += styled.String()
		index += ByteIndex(size)
		content = content[size:]

		runes = append(runes, VisualRune{
			C:     rune,
			Style: styled,
			VC:    styled.String(),
			Dirty: false,
		})
	}
	visualLines := strings.Split(util.Wrap(renderedLine, r.wrapWidth), "\n")

	// add empty line at the end of the paragraph with no line number
	visualLines = append(visualLines, "\n")
	runes = append(runes, VisualRune{C: '\n', Style: DefaultStyle, VC: "\n"})

	ret := []VisualLine{}
	start := 0
	for i, vl := range visualLines {
		ls := emptyLinum
		if i == 0 {
			ls = rlinum
		}
		stop := start + util.Len(vl)
		if len(runes) > 0 {
			runes[start].Dirty = true
		}
		ret = append(ret,
			VisualLine{
				BufferLinum: linum,
				Content:     vl,
				Runes:       runes[start:stop],
				// The wrap may cause the style left at the end of last line, then the linum style will cancel
				// the style of the first character in this line which will cause it's style to be lost.
				// But make the whole line dirty will cause the whole line to be rendered again which is
				// a waste of time. This could be solved by either rewriting the Wrap method or adding a
				// dirty flag to the VisualRune struct.
				Dirty:      true,
				LineNum:    ls,
				LinumStyle: linumStyle,
			},
		)
		start = stop
	}
	return ret
}

func (r *Renderer) GetBuffer() *Buffer {
	return r.buffer
}

func (r *Renderer) GetBufferX(line string, vy VisualLineIndex, vx VisualIndex) RuneIndex {
	return RuneIndex(util.LocBeforeWraped(line, r.wrapWidth, int(vx), int(vy)))
}

func (r *Renderer) GetVisualLineNumById(id epub.ManifestId) VisualLineIndex {
	return r.buffer.GetVisualLineNumById(id)
}

func (r *Renderer) MarkPosition(vy VisualLineIndex, vx VisualIndex) string {
	log.Debugf("MarkPosition: %d, %d", vy, vx)
	// TODO: show one space if the cursor is at an empty line
	vy = VisualLineIndex(util.MinInt(int(vy), len(r.buffer.visualLines)))
	if len(r.buffer.visualLines[vy].Content) > 0 {
		return r.buffer.visualLines[vy].MarkPosition(vx)
	}
	return ""
}

func (r *Renderer) ClearCursorStyles(vy VisualLineIndex) {
	r.buffer.visualLines[vy].ClearLine()
}

func (r *Renderer) MarkInline(vy VisualLineIndex, vxs, vxe VisualIndex) string {
	return r.buffer.visualLines[vy].MarkInline(vxs, vxe)
}

func (r *Renderer) MarkLine(vy VisualLineIndex) string {
	return r.buffer.visualLines[vy].MarkLine()
}
