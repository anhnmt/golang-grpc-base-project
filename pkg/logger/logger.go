package logger

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/natefinch/lumberjack/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewLogger the default logger
func NewLogger(logPath string) {
	// UNIX Time is faster and smaller than most timestamps
	consoleWriter := &zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	// Multi Writer
	writer := []io.Writer{
		getLogWriter(logPath),
		consoleWriter,
	}

	// Caller Marshal Function
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	log.Logger = zerolog.
		New(zerolog.MultiLevelWriter(writer...)).
		With().
		Timestamp().
		Caller().
		Logger()
}

// getLogWriter returns a lumberjack.logger
func getLogWriter(logFileUrl string) *lumberjack.Roller {
	options := &lumberjack.Options{
		MaxBackups: 5,  // Files
		MaxAge:     30, // 30 days
		Compress:   false,
	}

	// get log file path
	if logFileUrl == "" {
		logFileUrl = "logs/data.log"
	}

	var maxSize int64 = 100 * 1024 * 1024 // 100 MB
	roller, err := lumberjack.NewRoller(logFileUrl, maxSize, options)

	if err != nil {
		panic(err)
	}

	return roller
}
