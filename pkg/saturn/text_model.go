package saturn

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elinx/saturn/pkg/epub"
	"github.com/elinx/saturn/pkg/util"
	"github.com/elinx/saturn/pkg/viewport"
	log "github.com/sirupsen/logrus"
)

type textModel struct {
	book          *epub.Epub
	renderer      *Renderer
	prevModel     tea.Model
	viewport      viewport.Model
	width         int
	height        int
	currSectionId epub.ManifestId

	selectionStart Pos
	selectionEnd   Pos
	cursorReleased bool
}

func NewTextModel(book *epub.Epub, renderer *Renderer,
	currentId epub.ManifestId, prev tea.Model, width, height int) tea.Model {
	return &textModel{
		book:           book,
		renderer:       renderer,
		prevModel:      prev,
		width:          width,
		height:         height,
		currSectionId:  currentId,
		selectionStart: InvalidPos,
		selectionEnd:   InvalidPos,
		cursorReleased: true,
	}
}

func (m *textModel) Init() tea.Cmd {
	m.renderer.Render(m.width)
	m.viewport = viewport.New(m.width, m.height, m.renderer.buffer)
	m.viewport.Style = lipgloss.NewStyle()
	return nil
}

func (m *textModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m.prevModel, nil
		}
	case BlockMessage:
		id := message.(BlockMessage).ID
		pos := m.renderer.GetVisualLineNumById(id)
		m.viewport.SetYOffset(int(pos))
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			curr := Pos{
				X: util.MaxInt(0, msg.X-1),
				Y: msg.Y,
			}
			log.Debugf("mouse left clicked: %v", curr)
			if m.cursorReleased {
				m.clearCursor(m.selectionStart, m.selectionEnd)
				m.selectionStart = curr
			}
			m.clearCursor(m.selectionStart, m.selectionEnd)
			m.markSelection(m.selectionStart, curr)
			m.selectionEnd = curr
			m.cursorReleased = false
		case tea.MouseRelease:
			m.cursorReleased = true
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(message)
	return m, cmd
}

func (m *textModel) View() string {
	return m.viewport.View()
}

func (m *textModel) markPosition(p Pos) {
	visualLineNum := VisualLineIndex(p.Y + m.viewport.YOffset)
	m.renderer.MarkPosition(visualLineNum, VisualIndex(p.X))
}

func (m *textModel) markInline(sx, ex, sy int) {
	visualLineNum := VisualLineIndex(sy + m.viewport.YOffset)
	m.renderer.MarkInline(visualLineNum, VisualIndex(sx), VisualIndex(ex))
}

// sy should be smaller than ey
func (m *textModel) markCrossLine(sx, sy, ex, ey int) {
	m.markInline(sx, m.width-1, sy)
	for y := sy + 1; y < ey; y++ {
		m.markInline(0, m.width-1, y)
	}
	m.markInline(0, ex, ey)
}

func (m *textModel) markSelection(start, end Pos) {
	if start == end {
		m.markPosition(start)
		return
	}
	if start.Y == end.Y {
		sx := util.MinInt(start.X, end.X)
		ex := util.MaxInt(start.X, end.X)
		m.markInline(sx, ex, start.Y)
		return
	}
	if start.Y > end.Y {
		start, end = end, start
	}
	sx, sy := start.X, start.Y
	ex, ey := end.X, end.Y
	m.markCrossLine(sx, sy, ex, ey)
}

func (m *textModel) clearCursor(start, end Pos) {
	if start == InvalidPos || end == InvalidPos {
		return
	}
	if start.Y > end.Y {
		start, end = end, start
	}
	sy, ey := start.Y, end.Y
	for y := sy; y <= ey; y++ {
		visualLineNum := VisualLineIndex(y + m.viewport.YOffset)
		m.renderer.ClearCursorStyles(visualLineNum)
	}
}
