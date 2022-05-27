package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
	_ "github.com/elinx/saturn/pkg/logconfig"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("start app...")
	book := epub.NewEpub(os.Args[1])
	if err := book.Open(); err != nil {
		log.Fatal(err)
	}
	defer book.Close()

	program := tea.NewProgram(NewModel(book), tea.WithAltScreen())
	if err := program.Start(); err != nil {
		panic(err)
	}
}
