package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack/v3"
)

// New the default logger
func New(logFile string) {
	// Multi Writer
	writer := []io.Writer{
		os.Stdout,
	}

	if logFile != "" {
		roller, err := getLogWriter(logFile)
		if err != nil {
			panic(fmt.Errorf("get current directory failed: %s", err))
		}

		writer = append(writer, roller)
	}

	replacer := func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}

		return a
	}

	loggingLevel := new(slog.LevelVar)
	mutiWriter := io.MultiWriter(writer...)
	textHandler := slog.NewTextHandler(mutiWriter, &slog.HandlerOptions{
		AddSource:   true,
		Level:       loggingLevel,
		ReplaceAttr: replacer,
	})

	l := slog.New(textHandler)
	slog.SetDefault(l)
}

// getLogWriter returns a lumberjack.logger
func getLogWriter(logFile string) (*lumberjack.Roller, error) {
	var maxSize int64 = 50 * 1024 * 1024 // 50 MB

	options := &lumberjack.Options{
		MaxBackups: 5,  // files
		MaxAge:     30, // days
		Compress:   false,
	}

	roller, err := lumberjack.NewRoller(logFile, maxSize, options)
	if err != nil {
		return nil, err
	}

	return roller, nil
}
