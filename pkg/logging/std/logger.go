package std

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/protomem/gotube/pkg/logging"
)

var _ logging.Logger = (*Logger)(nil)

type Logger struct {
	*slog.Logger
	ctx context.Context
}

func New(level string, out io.Writer) (*Logger, error) {
	const op = "logger.New"

	parsedLevel, err := ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Logger{
		Logger: slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: parsedLevel})),
		ctx:    context.Background(),
	}, nil
}

func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.DebugContext(l.ctx, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.WarnContext(l.ctx, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.Logger.InfoContext(l.ctx, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.Logger.ErrorContext(l.ctx, msg, args...)
}

func (l *Logger) With(args ...any) logging.Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
		ctx:    l.ctx,
	}
}

func (l *Logger) Context() context.Context {
	return l.ctx
}

func (l *Logger) WithContext(ctx context.Context) logging.Logger {
	return &Logger{
		Logger: l.Logger,
		ctx:    ctx,
	}
}

func (l *Logger) Sync() error {
	return nil
}

func ParseLevel(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, errors.New("invalid log level")
	}
}
