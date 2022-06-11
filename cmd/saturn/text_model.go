package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elinx/saturn/pkg/epub"
	"github.com/elinx/saturn/pkg/parser"
	"github.com/muesli/reflow/wrap"
	log "github.com/sirupsen/logrus"
)

type textModel struct {
	book          *epub.Epub
	renderer      *parser.Renderer
	prevModel     tea.Model
	viewport      viewport.Model
	width         int
	height        int
	currSectionId epub.ManifestId

	selectionStart pos
	selectionEnd   pos
}

type pos struct {
	x, y int
}

var invalidPos pos = pos{-1, -1}

func NewTextModel(book *epub.Epub, renderer *parser.Renderer,
	currentId epub.ManifestId, prev tea.Model, width, height int) tea.Model {
	return &textModel{
		book:           book,
		renderer:       renderer,
		prevModel:      prev,
		viewport:       viewport.New(width, height),
		width:          width,
		height:         height,
		currSectionId:  currentId,
		selectionStart: invalidPos,
		selectionEnd:   invalidPos,
	}
}

func (m *textModel) Init() tea.Cmd {
	// content, err := m.book.GetContentByManifestId(m.currSectionId)
	// if err != nil {
	// 	log.Errorf("get content %v error: %v", m.currSectionId, err)
	// 	return tea.Quit
	// }
	// content = m.renderContent(content)
	content := m.renderer.Render(m.width)
	m.viewport.Style = lipgloss.NewStyle()
	// Bold(true).
	// Foreground(lipgloss.Color("#FAFAFA")).
	// Background(lipgloss.Color("#7D56F4")).
	// PaddingTop(2).
	// PaddingLeft(2)
	content = wrap.String(content, m.width)
	m.viewport.SetContent(content)
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
			// case "j":
			// 	m.continueNextPage()
			// case "k":
			// 	m.continuePrevPage()
		}
	case BlockMessage:
		id := message.(BlockMessage).ID
		pos := m.renderer.GetVisualPos(id)
		m.viewport.SetYOffset(pos)
	case tea.MouseMsg:
		switch msg.Type {
		// case tea.MouseWheelDown:
		// 	m.continueNextPage()
		// case tea.MouseWheelUp:
		// 	m.continuePrevPage()
		case tea.MouseLeft:
			log.Debugf("mouse left clicked: (%v, %v)", msg.X, msg.Y)
			if m.selectionStart == invalidPos {
				m.selectionStart = pos{msg.X, msg.Y}
			}
		case tea.MouseRelease:
			log.Debugf("mouse release: (%v, %v)", msg.X, msg.Y)
			m.selectionEnd = pos{msg.X, msg.Y}
			m.markSelection()
			m.selectionStart = invalidPos
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(message)
	return m, cmd
}

// func (m *textModel) renderContent(content string) string {
// if str, err := parser.Parse(content, parser.HtmlFormater{
// 	ColorProfile: termenv.ColorProfile(),
// 	Styles:       m.book.Styles,
// }); err != nil {
// 	log.Errorf("parse error: %v", err)
// 	return ""
// } else {
// 	return wrap.String(str, m.width)
// }
// 	return ""
// }

func (m *textModel) View() string {
	return m.viewport.View()
}

// func (m *textModel) continueNextPage() {
// 	if m.viewport.AtBottom() {
// 		if content, nextId, err := m.book.GetNextSection(m.currSectionId); err != nil {
// 			log.Info("get next section error: %v", err)
// 		} else {
// 			content = m.renderContent(content)
// 			m.content += content
// 			m.viewport.SetContent(m.content)
// 			m.currSectionId = nextId
// 		}
// 	}
// }

// func (m *textModel) continuePrevPage() {
// 	if m.viewport.AtTop() {
// 		if content, prevId, err := m.book.GetPrevSection(m.currSectionId); err != nil {
// 			log.Info("get prev section error: %v", err)
// 		} else {
// 			content = m.renderContent(content)
// 			m.content = content + m.content
// 			m.viewport.SetContent(m.content)
// 			m.currSectionId = prevId
// 		}
// 	}
// }

func (m *textModel) markSelection() {
	// if m.selectionStart == invalidPos || m.selectionEnd == invalidPos {
	// 	return
	// }
	// if m.selectionStart.y == m.selectionEnd.y {
	// 	m.viewport.Mark(m.selectionStart.y, m.selectionStart.x, m.selectionEnd.x)
	// } else {
	// 	m.viewport.Mark(m.selectionStart.y, m.selectionStart.x, -1)
	// 	m.viewport.Mark(m.selectionEnd.y, 0, m.selectionEnd.x)
	// }
}
