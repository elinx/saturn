package parser

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	cssparser "github.com/elinx/saturn/pkg/css_parser"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var DefaultFormater = HtmlFormater{
	ColorProfile: termenv.ColorProfile(),
}

type IHtmlFormater interface {
	Title(string) string
	I(string) string
	P(string, []html.Attribute) string
	PostProcess(string) string
}

type HtmlFormater struct {
	ColorProfile termenv.Profile
	Styles       []*cssparser.Rule
}

func (f HtmlFormater) Title(c string) string {
	return termenv.String(c).Bold().
		Foreground(f.ColorProfile.Color("#ffffff")).
		Background(f.ColorProfile.Color("#0000ff")).String()
}

func (f HtmlFormater) I(c string) string {
	return termenv.String(c).Italic().String()
}

func (f HtmlFormater) P(c string, attributes []html.Attribute) string {
	log.Println("attributes:", attributes)
	style := lipgloss.NewStyle()
	for _, attr := range attributes {
		if attr.Key == "class" {
			selector := "." + attr.Val
			for _, rule := range f.Styles {
				for _, s := range rule.Selector {
					if s == selector {
						log.Println("hit selector: ", selector)
						for _, prop := range rule.Declarations {
							if prop.Property == "color" {
								style.Foreground(lipgloss.Color(prop.Value))
							} else if prop.Property == "background-color" {
								style.Background(lipgloss.Color(prop.Value))
							} else if prop.Property == "border-top" {
								style.BorderTop(true)
							}
						}
					}
				}
			}
		}
	}
	rc := style.Render(c)
	log.Println("return:", rc)
	return rc
}

func (f HtmlFormater) PostProcess(c string) string {
	return strings.TrimSuffix(c, "\n")
}
