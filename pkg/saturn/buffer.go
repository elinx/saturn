package saturn

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/elinx/saturn/pkg/epub"
	"github.com/zyedidia/go-runewidth"
)

type IVisualVisiter interface {
	VisitBufferPrev(*Buffer) bool
	VisitBufferPost(*Buffer) bool
	VisitLinePrev(*VisualLine) bool
	VisitLinePost(*VisualLine) bool
	VisitRune(*VisualRune) bool
}

// RuneIndex returns the index of the rune in the given string
type RuneIndex int

// ByteIndex returns the index of the byte in the given string
type ByteIndex int

// VisualIndex returns the visual index of the screen position
type VisualIndex int

// BufferLineIndex is the index of the line in the buffer
type BufferLineIndex int

// VisualLineIndex is the index of the line in the rendered buffer
type VisualLineIndex int

// Segment is to describe the style in a part of the string
type Segment struct {
	Content string
	Style   string
	Pos     ByteIndex
}

// Line contains the text parsed from the ebooks together with
// it's styles specified. The Content field is the original text
// without rendering
type Line struct {
	Content  string
	Segments []Segment
	Style    string
}

type VisualRune struct {
	C     rune
	Style lipgloss.Style
	VC    string
	Dirty bool
}

func (r *VisualRune) Accept(visitor IVisualVisiter) {
	visitor.VisitRune(r)
}

type VisualLine struct {
	// index in the original buffer, some visual lines
	// may mapping to the same buffer line because of the
	// line wrapping.
	BufferLinum BufferLineIndex

	// rendered content
	Content string
	Runes   []VisualRune

	// rendered line number
	LineNum    string
	LinumStyle lipgloss.Style

	Dirty bool
}

func (v *VisualLine) Accept(visitor IVisualVisiter) {
	dirty := false
	dirty = dirty || visitor.VisitLinePrev(v)
	for _, vr := range v.Runes {
		dirty = dirty || visitor.VisitRune(&vr)
	}
	dirty = dirty || visitor.VisitLinePost(v)
	v.Dirty = dirty
}

func (v *VisualLine) MarkPosition(vx VisualIndex) string {
	v.Dirty = true
	pos := 0
	for i, vr := range v.Runes {
		width := runewidth.RuneWidth(vr.C)
		if pos+width > int(vx) {
			v.Runes[i].Style.Reverse(true)
			return string(v.Runes[i].C)
		}
		pos += width
	}
	return ""
}

func (v *VisualLine) MarkInline(start, end VisualIndex) string {
	v.Dirty = true
	content := ""
	prev := VisualIndex(-1)
	sameRune := func(curr, prev int) bool {
		if prev == -1 {
			return false
		}
		pos := 0
		for _, vr := range v.Runes {
			width := runewidth.RuneWidth(vr.C)
			if curr >= pos && curr < pos+width &&
				prev >= pos && prev < pos+width {
				return true
			}
			if pos > curr {
				return false
			}
			pos += width
		}
		return false
	}
	for x := start; x <= end; x++ {
		s := v.MarkPosition(x)
		if !sameRune(int(x), int(prev)) {
			content += s
		}
		prev = x
	}
	return content
}

func (v *VisualLine) MarkLine() string {
	v.Dirty = true
	content := ""
	for i := range v.Runes {
		content += string(v.Runes[i].C)
		v.Runes[i].Style.Reverse(true)
	}
	return content
}

func (v *VisualLine) ClearLine() {
	v.Dirty = true
	for i := range v.Runes {
		v.Runes[i].Style.Reverse(false)
	}
}

func (v *VisualLine) Replace(content string) {
	v.Dirty = false
	v.Content = content
}

func (v *VisualLine) Render() string {
	var content strings.Builder
	// content.WriteString(v.LineNum)
	for _, vr := range v.Runes {
		if vr.Dirty {
			content.WriteString(vr.Style.String())
		} else {
			content.WriteString(vr.VC)
		}
	}
	return content.String()
}

// Buffer is the ebook one to one mapping
type Buffer struct {
	renderer *Renderer
	Lines    []Line

	// The position of each block of the spine in the Lines
	BlockPos map[epub.ManifestId]BufferLineIndex

	// lineYOffsets is the offset of each line in the buffer after
	// being rendered to the screen. It is used to calculate the
	// position of each rune in the line.
	visualLineOffset []VisualLineIndex

	// visualLines are lines that have been rendered.
	visualLines []VisualLine
}

func NewBuffer() *Buffer {
	return &Buffer{
		Lines:       []Line{},
		BlockPos:    make(map[epub.ManifestId]BufferLineIndex),
		visualLines: make([]VisualLine, 0),
	}
}

// VisualLinesNum returns total lines number after rendition
func (b *Buffer) VisualLinesNum() int {
	return len(b.visualLines)
}

func (b *Buffer) LinesNum() int {
	return len(b.Lines)
}

// VisualLines return a portion of visual lines in a range
func (b *Buffer) VisualLines(start, end int) []string {
	return b.getVisualLines(VisualLineIndex(start), VisualLineIndex(end))
}

func (b *Buffer) getVisualLines(start, end VisualLineIndex) (res []string) {
	for i, line := range b.visualLines[start:end] {
		linum := VisualLineIndex(i) + start
		if line.Dirty {
			content := line.Render()
			b.visualLines[linum].Replace(content)
		}
		res = append(res, b.visualLines[linum].LineNum+b.visualLines[linum].Content)
	}
	return
}

func (b *Buffer) GetBufferLineNumById(id epub.ManifestId) BufferLineIndex {
	return b.BlockPos[id]
}

func (b *Buffer) GetVisualLineNumById(id epub.ManifestId) VisualLineIndex {
	return b.visualLineOffset[b.BlockPos[id]]
}

// GetBaseVisualLine returns the y position of the first line of the given
// visual index(one buffer line maybe rendered to multiple screen lines)
func (b *Buffer) GetBaseVisualLine(vy VisualLineIndex) VisualLineIndex {
	return b.visualLineOffset[b.GetBufferLineNumByVisual(vy)]
}

func (b *Buffer) GetBufferX(bufferLineNum BufferLineIndex, vy VisualLineIndex, vx VisualIndex) RuneIndex {
	vyBase := b.GetBaseVisualLine(vy)
	line := b.Lines[bufferLineNum].Content
	return b.renderer.GetBufferX(line, vy-vyBase, vx)
}

func (b *Buffer) GetBufferLineNumByVisual(visualLineNum VisualLineIndex) BufferLineIndex {
	for i, v := range b.visualLineOffset {
		if v == visualLineNum {
			return BufferLineIndex(i)
		} else if v > visualLineNum {
			return BufferLineIndex(i - 1)
		}
	}
	return BufferLineIndex(len(b.visualLineOffset) - 1)
}

type Visitor struct {
	Style   lipgloss.Style
	Content string
	Start   int
	End     int
}

func (v *Visitor) VisitBufferPrev(b *Buffer) bool {
	return false
}

func (v *Visitor) VisitBufferPost(b *Buffer) bool {
	return false
}

func (v *Visitor) VisitLinePrev(l *VisualLine) bool {
	return false
}

func (v *Visitor) VisitLinePost(l *VisualLine) bool {
	return false
}

func (v *Visitor) VisitRune(r *VisualRune) bool {
	v.Content += string(r.C)
	r.Style = r.Style.Reverse(true)
	return true
}

func (b *Buffer) Accept(visitor IVisualVisiter) {
	for _, line := range b.visualLines {
		line.Accept(visitor)
	}
}

func visit(b *Buffer) {
	visitor := &Visitor{}
	b.Accept(visitor)
}
