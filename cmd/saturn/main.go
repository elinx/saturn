package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
	"golang.org/x/net/html"
)

func parse(content string) (string, error) {
	htmlNode, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", err
	}
	var result string
	for c := htmlNode.FirstChild; c != nil; c = c.NextSibling {
		fmt.Printf("%+v -----\n", c)
		if c.Type == html.ElementNode {
			if c.Data == "body" {
				for c := c.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						result += c.Data
					}
				}
			}
		}
	}
	return result, nil
}

func main() {
	book := epub.NewEpub("test/data/TaoTeChing.epub")
	book.Open()
	defer book.Close()

	// for i, v := range book.Toc.NavMap.NavPoints {
	// 	fmt.Printf("%d \n", i)
	// 	if len(v.NavPoints) > 0 {
	// 		for _, v := range v.NavPoints {
	// 			fmt.Printf("\t%+v\n", v)
	// 		}
	// 	}
	// }
	// contentHtml, err := book.GetTableOfContent()
	// if err != nil {
	// 	panic(err)
	// }
	// // fmt.Println(contentHtml)
	// content, err := parse(contentHtml)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(content)
	program := tea.NewProgram(NewModel(book), tea.WithAltScreen())
	if err := program.Start(); err != nil {
		panic(err)
	}
}
