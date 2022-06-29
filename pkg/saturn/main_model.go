package saturn

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
	log "github.com/sirupsen/logrus"
)

type item struct {
	title string
	src   epub.ManifestId
}

func (i item) FilterValue() string  { return i.title }
func (i item) Title() string        { return i.title }
func (i item) Description() string  { return "" }
func (i item) Src() epub.ManifestId { return i.src }

func newItems(book *epub.Epub) []list.Item {
	content := []list.Item{}
	toc := book.GetTableOfContent()
	for _, v := range toc.Orders {
		record := toc.Content[v]
		content = append(content, item{
			title: v,
			src:   record.ID,
		})
	}
	return content
}

func NewMainModel(book *epub.Epub, renderer *Renderer) tea.Model {
	return &mainModel{
		book:     book,
		renderer: renderer,
		tocModel: list.New(newItems(book), list.DefaultDelegate{
			ShowDescription: false,
			Styles:          list.NewDefaultItemStyles(),
		}, 30, 30),
	}
}

type mainModel struct {
	book      *epub.Epub
	renderer  *Renderer
	tocModel  list.Model
	textModel tea.Model
	width     int
	height    int
}

func (m *mainModel) Init() tea.Cmd {
	return nil
}

type BlockMessage struct {
	ID  epub.ManifestId
	Msg string
}

func (m *mainModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if item, ok := m.tocModel.SelectedItem().(item); !ok {
				return m, tea.Quit
			} else {
				log.Debugf("item selected: %s", item.Src())
				return m.textModel, func() tea.Msg {
					return BlockMessage{item.Src(), "select from toc"}
				}
			}
		}
	case tea.WindowSizeMsg:
		log.Debugf("window size changed: ", msg.Width, msg.Height)
		m.width = msg.Width
		m.height = msg.Height
		m.textModel = NewTextModel(m.book, m.renderer,
			m.tocModel.SelectedItem().(item).Src(), m, m.width, m.height)
		m.textModel.Init()
	}

	var cmd tea.Cmd
	m.tocModel, cmd = m.tocModel.Update(message)
	return m, cmd
}

func (m *mainModel) View() string {
	return m.tocModel.View()
}
