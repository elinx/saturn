package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
	"github.com/elinx/saturn/pkg/parser"
	"github.com/muesli/reflow/wrap"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
)

type textModel struct {
	book          *epub.Epub
	prevModel     tea.Model
	content       string
	viewport      viewport.Model
	width         int
	height        int
	currSectionId epub.ManifestId
}

func NewTextModel(book *epub.Epub, href epub.HRef, prev tea.Model, width, height int) tea.Model {
	return &textModel{
		book:          book,
		prevModel:     prev,
		viewport:      viewport.New(width, height),
		width:         width,
		height:        height,
		currSectionId: book.HrefToManifestId(href),
	}
}

func (m *textModel) Init() tea.Cmd {
	content, err := m.book.GetContentByManifestId(m.currSectionId)
	if err != nil {
		log.Errorf("get content %v error: %v", m.currSectionId, err)
		return tea.Quit
	}
	content = m.renderContent(content)
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
		case "j":
			m.continueNextPage()
		case "k":
			m.continuePrevPage()
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelDown:
			m.continueNextPage()
		case tea.MouseWheelUp:
			m.continuePrevPage()
		case tea.MouseLeft:
			log.Infof("mouse left clicked: (%v, %v)", msg.X, msg.Y)
		case tea.MouseRelease:
			log.Infof("mouse release: (%v, %v)", msg.X, msg.Y)
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(message)
	return m, cmd
}

func (m *textModel) renderContent(content string) string {
	if str, err := parser.Parse(content, parser.HtmlFormater{
		ColorProfile: termenv.ColorProfile(),
		Styles:       m.book.Styles,
	}); err != nil {
		log.Errorf("parse error: %v", err)
		return ""
	} else {
		return wrap.String(str, m.width)
	}
}

func (m *textModel) View() string {
	return m.viewport.View()
}

func (m *textModel) continueNextPage() {
	if m.viewport.AtBottom() {
		if content, nextId, err := m.book.GetNextSection(m.currSectionId); err != nil {
			log.Info("get next section error: %v", err)
		} else {
			content = m.renderContent(content)
			m.content += content
			m.viewport.SetContent(m.content)
			m.currSectionId = nextId
		}
	}
}

func (m *textModel) continuePrevPage() {
	if m.viewport.AtTop() {
		if content, prevId, err := m.book.GetPrevSection(m.currSectionId); err != nil {
			log.Info("get prev section error: %v", err)
		} else {
			content = m.renderContent(content)
			m.content = content + m.content
			m.viewport.SetContent(m.content)
			m.currSectionId = prevId
		}
	}
}
