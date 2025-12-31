package logging

import (
	"io"
	"log/slog"
	"os"
)

type ComiketLogger struct {
	file *os.File
	*slog.Logger
}

var Logger ComiketLogger

func InitializeLogging(rawLogLevel string, filePath string) error {
	var w io.Writer
	var f *os.File
	if filePath != "" {
		var err error
		f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		w = io.MultiWriter(os.Stdout, f)
	} else {
		w = os.Stdout
	}

	var logLevel slog.Level
	err := logLevel.UnmarshalText([]byte(rawLogLevel))
	if err != nil {
		return err
	}

	// TODO: logging options?
	th := slog.NewTextHandler(w, nil)
	Logger = ComiketLogger{f, slog.New(NewLevelHandler(logLevel, th))}
	return nil

}

func (comiketLogger *ComiketLogger) TeardownLogging() {
	if comiketLogger.file != nil {
		comiketLogger.file.Close()
	}
}
