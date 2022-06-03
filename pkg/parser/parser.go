package parser

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

func parse(node *html.Node, formater IHtmlFormater) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	var result string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		result += parse(c, formater)
	}
	if node.Type == html.ElementNode {
		switch node.Data {
		case "i":
			return formater.I(result)
		case "p":
			return formater.P(result, node.Attr)
		case "title":
			return formater.Title(result)
		case "h1", "h2", "h3", "h4", "h5", "h6":
			return formater.Header(result)
		}
	}
	return result
}

func Parse(content string, formater IHtmlFormater) (string, error) {
	log.Infoln("Enter into parsing of HTML")
	htmlNode, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", err
	}
	return formater.PostProcess(parse(htmlNode, formater)), nil
}
