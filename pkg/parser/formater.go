package parser

import (
	"strings"

	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var DefaultFormater = htmlFormater{
	colorProfile: termenv.ColorProfile(),
}

type IHtmlFormater interface {
	Title(string) string
	I(string) string
	P(string, []html.Attribute) string
	PostProcess(string) string
}

type htmlFormater struct {
	colorProfile termenv.Profile
}

func (f htmlFormater) Title(c string) string {
	return termenv.String(c).Bold().
		Foreground(f.colorProfile.Color("#ffffff")).
		Background(f.colorProfile.Color("#0000ff")).String()
}

func (f htmlFormater) I(c string) string {
	return termenv.String(c).Italic().String()
}

func (f htmlFormater) P(c string, attributes []html.Attribute) string {
	log.Infof("attributes: %v\n", attributes)
	return c
}

func (f htmlFormater) PostProcess(c string) string {
	return strings.TrimSuffix(c, "\n")
}
