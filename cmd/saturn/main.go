package main

import (
	"fmt"

	"github.com/elinx/saturn/pkg/epub"
)

func main() {
	book := epub.NewEpub("test/data/TaoTeChing.epub")
	book.OpenFile()
	defer book.Close()
	fmt.Println(book.GetChapterByIndex(0))
}
