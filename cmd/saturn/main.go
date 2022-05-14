package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elinx/saturn/pkg/epub"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

const (
	logFilename = "log.txt"
)

type simpleFormatter struct {
	log.TextFormatter
}

func (f *simpleFormatter) Format(entry *log.Entry) ([]byte, error) {
	// this whole mess of dealing with ansi color codes is required if you want the colored output otherwise you will lose colors in the log levels
	var levelColor int
	switch entry.Level {
	case log.DebugLevel, log.TraceLevel:
		levelColor = 31 // gray
	case log.WarnLevel:
		levelColor = 33 // yellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}
	return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m %s %s:%d %s\n",
		levelColor, strings.ToUpper(entry.Level.String()),
		entry.Time.Format(f.TimestampFormat),
		path.Base(entry.Caller.File), entry.Caller.Line,
		entry.Message)), nil
}

func init() {
	logFile, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(logFile)
	log.SetReportCaller(true)
	log.SetFormatter(&simpleFormatter{
		log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	})
}

func parse(content string) (string, error) {
	htmlNode, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", err
	}
	var result string
	for c := htmlNode.FirstChild; c != nil; c = c.NextSibling {
		log.Printf("%+v -----\n", c)
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
