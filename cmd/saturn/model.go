package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
)

func NewModel(book *epub.Epub) tea.Model {
	return model{
		book: book,
	}
}

type model struct {
	book *epub.Epub
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	content := ""
	for _, v := range m.book.Toc.NavMap.NavPoints {
		content += v.NavLable.Text + "\n"
		for _, v := range v.NavPoints {
			content += "  " + v.NavLable.Text + "\n"
		}
	}
	return content
}
