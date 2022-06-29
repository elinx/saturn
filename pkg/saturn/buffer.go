package saturn

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/elinx/saturn/pkg/epub"
	"github.com/zyedidia/go-runewidth"
)

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
}

type VisualLine struct {
	// index in in the original buffer, some visual lines
	// may mapping to the same buffer line because of the
	// line wrapping.
	BufferLinum BufferLineIndex

	// rendered content
	Content string
	Runes   []VisualRune

	Dirty bool
}

func (v *VisualLine) MarkPosition(vx VisualIndex) {
	v.Dirty = true
	pos := 0
	for i, vr := range v.Runes {
		width := runewidth.RuneWidth(vr.C)
		if pos+width > int(vx) {
			v.Runes[i].Style.Reverse(true)
			return
		}
		pos += width
	}
}

func (v *VisualLine) MarkLine() {
	v.Dirty = true
	for i := range v.Runes {
		v.Runes[i].Style.Reverse(true)
	}
}

func (v *VisualLine) ClearLine() {
	v.Dirty = true
	for i := range v.Runes {
		v.Runes[i].Style.Reverse(false)
	}
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

// VisualLines return a portion of visual lines in a range
func (b *Buffer) VisualLines(start, end int) []string {
	return b.getVisualLines(VisualLineIndex(start), VisualLineIndex(end))
}

func (b *Buffer) getVisualLines(start, end VisualLineIndex) (res []string) {
	for i, line := range b.visualLines[start:end] {
		linum := VisualLineIndex(i) + start
		content := ""
		if line.Dirty {
			for _, vr := range line.Runes {
				content += vr.Style.String()
			}
			b.visualLines[linum].Content = content
			b.visualLines[linum].Dirty = false
		} else {
			content = line.Content
		}
		res = append(res, content)
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
