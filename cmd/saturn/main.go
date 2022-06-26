package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
	_ "github.com/elinx/saturn/pkg/logconfig"
	"github.com/elinx/saturn/pkg/saturn"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("start app...")
	book := epub.NewEpub(os.Args[1])
	if err := book.Open(); err != nil {
		log.Fatal(err)
	}
	defer book.Close()

	parser := saturn.NewParser(book)
	if err := parser.Parse(); err != nil {
		log.Fatal(err)
	}
	renderer := saturn.NewRender(book, parser.GetBuffer())

	program := tea.NewProgram(saturn.NewMainModel(book, renderer),
		tea.WithAltScreen(), tea.WithMouseAllMotion())
	if err := program.Start(); err != nil {
		panic(err)
	}
}
