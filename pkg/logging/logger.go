package logging

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type Logger struct {
	*slog.Logger
}

func New(out io.Writer, lvlStr string) (*Logger, error) {
	lvl, err := ParseLevel(lvlStr)
	if err != nil {
		return nil, fmt.Errorf("logging.New: %w", err)
	}

	return &Logger{slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: lvl}))}, nil
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{l.Logger.With(args...)}
}

func ParseLevel(lvlStr string) (slog.Level, error) {
	switch strings.ToUpper(lvlStr) {
	case slog.LevelDebug.String():
		return slog.LevelDebug, nil
	case slog.LevelInfo.String():
		return slog.LevelInfo, nil
	case slog.LevelWarn.String():
		return slog.LevelWarn, nil
	case slog.LevelError.String():
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("logging.ParseLevel: %w", errors.New("log level unsupported"))
	}
}
