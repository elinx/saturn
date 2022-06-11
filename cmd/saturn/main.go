package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
	_ "github.com/elinx/saturn/pkg/logconfig"
	"github.com/elinx/saturn/pkg/parser"
	log "github.com/sirupsen/logrus"
	tcell "github.com/zyedidia/tcell/v2"
)

type TextCell struct {
	c     rune
	X, Y  int
	style tcell.Style
}
type mouseModel struct {
	content        string
	texts          []TextCell
	selectionStart pos
	selectionEnd   pos
}

func (m *mouseModel) Init() tea.Cmd {
	for i, c := range m.content {
		m.texts = append(m.texts, TextCell{c, i, 0, tcell.StyleDefault})
	}
	tcell.NewScreen()
	return nil
}

func (m *mouseModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			log.Debugf("mouse left clicked: (%v, %v)", msg.X, msg.Y)
			if m.selectionStart == invalidPos {
				m.selectionStart = pos{msg.X, msg.Y}
			}
			m.texts[msg.X].style = m.texts[msg.X].style.Reverse(true)
		case tea.MouseRelease:
			log.Debugf("mouse release: (%v, %v)", msg.X, msg.Y)
			m.selectionEnd = pos{msg.X, msg.Y}
			// m.markSelection()
			m.selectionStart = invalidPos
		}
	}
	return m, nil
}

func (m *mouseModel) View() string {
	res := ""
	for _, t := range m.texts {
		_, _, attrs := t.style.Decompose()
		if attrs&tcell.AttrReverse != 0 {
			res += "\x1b[7m"
			res += string(t.c)
			res += "\x1b[0m"
		} else {
			res += string(t.c)
		}
	}
	return res
}

func (m *mouseModel) markSelection() {
	if m.selectionStart == invalidPos || m.selectionEnd == invalidPos {
		return
	}
	if m.selectionEnd.y == m.selectionStart.y {
		for x := m.selectionStart.x; x <= m.selectionEnd.x; x++ {
			m.texts[x].style = m.texts[x].style.Reverse(true)
		}
	}
}

func NewMouseModel() tea.Model {
	return &mouseModel{
		content:        "mouse model example",
		selectionStart: invalidPos,
		selectionEnd:   invalidPos,
	}
}

func main() {
	log.Println("start app...")
	book := epub.NewEpub(os.Args[1])
	if err := book.Open(); err != nil {
		log.Fatal(err)
	}
	defer book.Close()

	renderer := parser.New(book)
	if err := renderer.Parse(); err != nil {
		log.Fatal(err)
	}

	program := tea.NewProgram(NewModel(book, renderer),
		tea.WithAltScreen(), tea.WithMouseAllMotion())
	if err := program.Start(); err != nil {
		panic(err)
	}
}
