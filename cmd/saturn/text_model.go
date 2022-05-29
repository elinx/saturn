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
	book      *epub.Epub
	file      string
	prevModel tea.Model
	content   string
	viewport  viewport.Model
}

func NewTextModel(book *epub.Epub, id string, prev tea.Model, width, height int) tea.Model {
	file := book.GetFullPath(id)
	content := readContent(book, file)
	content = wrap.String(content, width)
	viewport := viewport.New(width, height)
	viewport.SetContent(content)
	return textModel{
		book:      book,
		file:      file,
		prevModel: prev,
		content:   content,
		viewport:  viewport,
	}
}

func (m textModel) Init() tea.Cmd {
	log.Println("content: ", m.content)
	return nil
}

func (m textModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m.prevModel, nil
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(message)
	return m, cmd
}

func readContent(book *epub.Epub, file string) string {
	if content, err := book.GetContentByFilePath(file); err != nil {
		return err.Error()
	} else {
		if str, err := parser.Parse(content, parser.HtmlFormater{
			ColorProfile: termenv.ColorProfile(),
			Styles:       book.Styles,
		}); err != nil {
			log.Errorf("parse error: %v", err)
			return err.Error()
		} else {
			return str
		}
	}
}

func (m textModel) View() string {
	return m.viewport.View()
}
