package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	cssparser "github.com/elinx/saturn/pkg/css_parser"
	"github.com/elinx/saturn/pkg/epub"
	_ "github.com/elinx/saturn/pkg/logconfig"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("start app...")
	book := epub.NewEpub(os.Args[1])
	book.Open()
	defer book.Close()

	for i, v := range book.GetCssFiles() {
		log.Printf("%d: %s", i, v)
		if content, err := book.GetContentByFilePath(v); err == nil {
			log.Println(cssparser.NewParser().Parse(content))
		}
	}

	// for i, v := range book.Toc.NavMap.NavPoints {
	// 	log.Printf("%d \n", i)
	// 	if len(v.NavPoints) > 0 {
	// 		for _, v := range v.NavPoints {
	// 			log.Printf("\t%+v\n", v)
	// 		}
	// 	}
	// }
	// contentHtml, err := book.GetTableOfContent()
	// if err != nil {
	// 	panic(err)
	// }
	// // log.Println(contentHtml)
	// content, err := parse(contentHtml)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(content)
	program := tea.NewProgram(NewModel(book), tea.WithAltScreen())
	if err := program.Start(); err != nil {
		panic(err)
	}
}
