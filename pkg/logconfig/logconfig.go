package logconfig

import (
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
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
	logFile, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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
