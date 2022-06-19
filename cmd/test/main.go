package main

import (
	"fmt"
	"strings"

	"github.com/muesli/reflow/wrap"
)

func render(line string) string {
	// `ghi jk` are marked as italic: \x1b[3m
	lineStyle := "ab c d ef \x1b[3mghi jk\x1b[0m lmn opqrst uvw x yz\n"
	lineWraped := wrap.String(lineStyle, 10)
	return lineWraped
}

func main() {
	line := "ab c d ef ghi jk lmn opqrst uvw x yz"
	lineRendered := render(line)
	fmt.Println(line, string(line[28]))
	fmt.Println("---")
	lines := strings.Split(lineRendered, "\n")
	fmt.Print(lineRendered, string(lines[2][7]))

	fmt.Println("---")
	fmt.Println(wrap.String("iamamedicoreprogrammer", 5))
}
