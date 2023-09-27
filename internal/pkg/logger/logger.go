package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bytedance/sonic"
	"github.com/natefinch/lumberjack/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New the default logger
func New(logFile string) {
	// UNIX Time is faster and smaller than most timestamps
	consoleWriter := &zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	// Multi Writer
	writer := []io.Writer{
		consoleWriter,
	}

	if logFile != "" {
		roller, err := getLogWriter(logFile)
		if err != nil {
			log.Panic().Msg(fmt.Sprintf("get current directory failed: %v", err))
		}

		writer = append(writer, roller)
	}

	// Caller Marshal Function
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	zerolog.InterfaceMarshalFunc = sonic.Marshal

	l := zerolog.
		New(zerolog.MultiLevelWriter(writer...)).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = l
	zerolog.DefaultContextLogger = &l
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
