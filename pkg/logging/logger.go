package logging

import "context"

type Logger interface {
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)

	With(args ...any) Logger

	Context() context.Context
	WithContext(ctx context.Context) Logger

	Extractor(fn func(context.Context) []any)

	Sync() error
}
