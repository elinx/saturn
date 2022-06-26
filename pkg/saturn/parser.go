package saturn

import (
	"strings"

	"github.com/elinx/saturn/pkg/epub"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type Parser struct {
	book   *epub.Epub
	buffer *Buffer
}

func NewParser(book *epub.Epub) *Parser {
	return &Parser{
		book:   book,
		buffer: NewBuffer(),
	}
}

func (p *Parser) GetBuffer() *Buffer {
	return p.buffer
}

// Parse iterates over the spine and parses each HTML file
func (p *Parser) Parse() error {
	content, err := p.book.GetSpinContent()
	if err != nil {
		return err
	}
	for _, id := range content.Orders {
		htmlContent := content.Contents[id]
		p.buffer.BlockPos[id] = BufferLineIndex(len(p.buffer.Lines))
		p.parse1(htmlContent)
	}
	return nil
}

func (p *Parser) parse1(content string) error {
	log.Infoln("Enter into parsing of HTML")
	htmlNode, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return err
	}
	if _, err := p.parse2(htmlNode); err != nil {
		return err
	}
	return nil
}

func (p *Parser) parse2(n *html.Node) (*Segment, error) {
	switch n.Type {
	case html.TextNode:
		if len(strings.TrimSpace(n.Data)) == 0 {
			return nil, nil
		}
		return &Segment{n.Data, "", 0}, nil
	case html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if _, err := p.parse2(c); err != nil {
				return nil, err
			}
		}
		return nil, nil
	case html.CommentNode:
		return nil, nil
	}
	var segments []Segment
	pos := ByteIndex(0)
	contents := []string{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		segment, err := p.parse2(c)
		if err != nil {
			return nil, err
		}
		if segment != nil {
			segment.Pos = pos
			segments = append(segments, *segment)
			contents = append(contents, segment.Content)
			pos += ByteIndex(len(segment.Content))
		}
	}
	lineContent := strings.Join(contents, "")
	switch n.Data {
	case "head", "html", "body", "link":
		// ignore
	case "svg", "image", "img":
		// TODO: support image display
	case "style":
		// TODO: support inline style
	case "i", "b", "strong", "span", "em":
		if n.Parent.Data == "body" {
			p.buffer.Lines = append(p.buffer.Lines, Line{lineContent, segments, n.Data})
			return nil, nil
		}
		return &Segment{lineContent, n.Data, 0}, nil
	default:
		p.buffer.Lines = append(p.buffer.Lines, Line{lineContent, segments, n.Data})
	}
	return nil, nil
}
