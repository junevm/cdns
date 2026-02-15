package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"cli/internal/config"

	"golang.org/x/term"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey string

const loggerKey contextKey = "logger"

// New creates a new slog.Logger based on configuration
// It automatically detects TTY and adjusts output format accordingly
func New(cfg config.LoggerConfig) *slog.Logger {
	var handler slog.Handler

	// Determine log level
	level := parseLevel(cfg.Level)
	levelVar := &slog.LevelVar{}
	levelVar.Set(level)

	// Check if output is a TTY
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))

	opts := &slog.HandlerOptions{
		Level:     levelVar,
		AddSource: level == slog.LevelDebug,
	}

	// Use JSON format for non-TTY (CI/CD) or if explicitly requested
	if !isTTY || cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Use pretty text handler for TTY (development)
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// SetLevel updates the logger level globally if it was created with LevelVar
func SetLevel(logger *slog.Logger, level slog.Level) {
	// This only works if we have a way to access the LevelVar.
	// Since slog doesn't expose it easily from the logger,
	// we might need a better way.
}

// parseLevel converts string level to slog.Level
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// WithLogger adds a logger to the context
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext extracts the logger from the context
// Returns a default logger if none is found
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// NewNoop creates a no-op logger (useful for testing)
func NewNoop() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
