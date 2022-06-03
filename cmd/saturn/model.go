package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
	log "github.com/sirupsen/logrus"
)

type item struct {
	title string
	src   epub.HRef
}

func (i item) FilterValue() string { return i.title }
func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) Src() epub.HRef      { return i.src }

func newItems(book *epub.Epub) []list.Item {
	content := []list.Item{}
	for _, v := range book.Toc.NavMap.NavPoints {
		content = append(content, item{title: v.NavLable.Text, src: v.Content.Src})
		for _, v := range v.NavPoints {
			content = append(content, item{title: "  " + v.NavLable.Text, src: v.Content.Src})
		}
	}
	return content
}

func NewModel(book *epub.Epub) tea.Model {
	return model{
		book: book,
		list: list.New(newItems(book), list.DefaultDelegate{
			ShowDescription: false,
			Styles:          list.NewDefaultItemStyles(),
		}, 30, 30),
	}
}

type model struct {
	book   *epub.Epub
	list   list.Model
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(item); !ok {
				return m, tea.Quit
			} else {
				log.Debugf("item selected: %s", item.Src())
				model := NewTextModel(m.book, item.Src(), m, m.width, m.height)
				return model, model.Init()
			}
		}
	case tea.WindowSizeMsg:
		log.Debugf("window size changed: ", msg.Width, msg.Height)
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(message)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}
