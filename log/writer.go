package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type LevelWriter struct {
	errWriter   io.Writer
	debugWriter io.Writer
}

func NewLevelWriter() LevelWriter {
	return LevelWriter{
		errWriter:   os.Stderr,
		debugWriter: os.Stdout,
	}
}

func NewConsoleLevelWriter() LevelWriter {
	return LevelWriter{
		errWriter:   zerolog.ConsoleWriter{Out: os.Stderr},
		debugWriter: zerolog.ConsoleWriter{Out: os.Stdout},
	}
}

func (w LevelWriter) Write(p []byte) (n int, err error) {
	return w.debugWriter.Write(p)
}
