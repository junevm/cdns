package logger

import (
	"log/slog"
	"testing"

	"github.com/junevm/cdns/internal/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config config.LoggerConfig
	}{
		{
			name: "info level text format",
			config: config.LoggerConfig{
				Level:  "info",
				Format: "text",
			},
		},
		{
			name: "debug level json format",
			config: config.LoggerConfig{
				Level:  "debug",
				Format: "json",
			},
		},
		{
			name: "error level",
			config: config.LoggerConfig{
				Level:  "error",
				Format: "text",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.config)
			if logger == nil {
				t.Fatal("expected logger to be created")
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"warning", slog.LevelWarn},
		{"error", slog.LevelError},
		{"invalid", slog.LevelInfo}, // defaults to info
		{"", slog.LevelInfo},        // defaults to info
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewNoop(t *testing.T) {
	logger := NewNoop()
	if logger == nil {
		t.Fatal("expected noop logger to be created")
	}

	// Should not panic when logging
	logger.Info("test message")
	logger.Error("test error")
}
