package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger is the structured logging interface.
type Logger interface {
	Infow(ctx context.Context, msg string, keysAndValues ...any)
	Warnw(ctx context.Context, msg string, keysAndValues ...any)
	Errorw(ctx context.Context, msg string, keysAndValues ...any)
}

type logger struct {
	internal *slog.Logger
}

// Errorw logs an error-level message with structured key-value pairs.
func (l *logger) Errorw(ctx context.Context, msg string, keysAndValues ...any) {
	l.internal.ErrorContext(ctx, msg, keysAndValues...)
}

// Infow logs an info-level message with structured key-value pairs.
func (l *logger) Infow(ctx context.Context, msg string, keysAndValues ...any) {
	l.internal.InfoContext(ctx, msg, keysAndValues...)
}

// NewLogger creates and returns a new Logger backed by log/slog.
func NewLogger() Logger {
	return &logger{internal: slog.New(slog.NewTextHandler(os.Stdout, nil))}
}

// Warnw logs a warn-level message with structured key-value pairs.
func (l *logger) Warnw(ctx context.Context, msg string, keysAndValues ...any) {
	l.internal.WarnContext(ctx, msg, keysAndValues...)
}
