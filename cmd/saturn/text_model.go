package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
	"github.com/elinx/saturn/pkg/parser"
)

type textModel struct {
	book      *epub.Epub
	file      string
	prevModel tea.Model
}

func NewTextModel(book *epub.Epub, id string, prev tea.Model) tea.Model {
	return textModel{
		book:      book,
		file:      book.GetFullPath(id),
		prevModel: prev,
	}
}

func (m textModel) Init() tea.Cmd {
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
	return m, nil
}

func (m textModel) View() string {
	if content, err := m.book.GetContentByFilePath(m.file); err != nil {
		return err.Error()
	} else {
		if str, err := parser.Parse(content, parser.DefaultFormater); err != nil {
			return err.Error()
		} else {
			return str
		}
	}
}
