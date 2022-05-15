package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
	_ "github.com/elinx/saturn/pkg/logconfig"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("start app...")
	book := epub.NewEpub("test/data/TaoTeChing.epub")
	book.Open()
	defer book.Close()

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
