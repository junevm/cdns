package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader()
	if loader == nil {
		t.Fatal("expected loader to be created")
	}
	if loader.k == nil {
		t.Fatal("expected koanf instance to be initialized")
	}
}

func TestLoadDefaults(t *testing.T) {
	loader := NewLoader()
	cfg, err := loader.Load("", nil)
	if err != nil {
		t.Fatalf("unexpected error loading defaults: %v", err)
	}

	if cfg.Logger.Level != "warn" {
		t.Errorf("expected logger level to be 'warn', got %s", cfg.Logger.Level)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	// Set environment variables
	_ = os.Setenv("APP_LOGGER_LEVEL", "debug")
	defer func() {
		_ = os.Unsetenv("APP_LOGGER_LEVEL")
	}()

	loader := NewLoader()
	cfg, err := loader.Load("", nil)
	if err != nil {
		t.Fatalf("unexpected error loading from environment: %v", err)
	}

	if cfg.Logger.Level != "debug" {
		t.Errorf("expected logger level to be 'debug', got %s", cfg.Logger.Level)
	}
}

func TestLoadFromFlags(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("logger.level", "", "")

	_ = flags.Set("logger.level", "warn")

	loader := NewLoader()
	cfg, err := loader.Load("", flags)
	if err != nil {
		t.Fatalf("unexpected error loading from flags: %v", err)
	}

	if cfg.Logger.Level != "warn" {
		t.Errorf("expected logger level to be 'warn', got %s", cfg.Logger.Level)
	}
}

func TestLoadPriority(t *testing.T) {
	// Set environment variable
	_ = os.Setenv("APP_LOGGER_LEVEL", "debug")
	defer func() { _ = os.Unsetenv("APP_LOGGER_LEVEL") }()

	// Set flag (should override environment)
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("logger.level", "", "")
	_ = flags.Set("logger.level", "error")

	loader := NewLoader()
	cfg, err := loader.Load("", flags)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	// Flag should win
	if cfg.Logger.Level != "error" {
		t.Errorf("expected logger level to be 'error' (from flag), got %s", cfg.Logger.Level)
	}
}
