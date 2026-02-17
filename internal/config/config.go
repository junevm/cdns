package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

// DefaultConfig is the default YAML configuration
const DefaultConfig = `
logger:
  level: warn
  format: text

dns:
  default_scope: active
  default_interfaces: []
  custom_presets:
    personal: ["1.1.1.1", "1.0.0.1"]
`

// Config represents the application configuration
type Config struct {
	Logger     LoggerConfig `koanf:"logger"`
	DNS        DNSConfig    `koanf:"dns"`
	LoadedFrom string       `koanf:"-"` // Not loaded from config, but set by loader
}

// Validate ensures the configuration is valid
func (c *Config) Validate() error {
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[strings.ToLower(c.Logger.Level)] {
		return fmt.Errorf("invalid logger.level: %s", c.Logger.Level)
	}

	validFormats := map[string]bool{"text": true, "json": true}
	if !validFormats[strings.ToLower(c.Logger.Format)] {
		return fmt.Errorf("invalid logger.format: %s", c.Logger.Format)
	}

	validScopes := map[string]bool{"active": true, "all": true, "explicit": true}
	if c.DNS.DefaultScope != "" && !validScopes[strings.ToLower(c.DNS.DefaultScope)] {
		return fmt.Errorf("invalid dns.default_scope: %s", c.DNS.DefaultScope)
	}

	return nil
}

// LoggerConfig contains logging settings
type LoggerConfig struct {
	Level  string `koanf:"level"`
	Format string `koanf:"format"` // "text" or "json"
}

// DNSConfig contains DNS-specific settings
type DNSConfig struct {
	DefaultScope      string              `koanf:"default_scope"`
	DefaultInterfaces []string            `koanf:"default_interfaces"`
	CustomPresets     map[string][]string `koanf:"custom_presets"`
}

// Loader handles configuration loading using Koanf
type Loader struct {
	k *koanf.Koanf
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	return &Loader{
		k: koanf.New("."),
	}
}

// DefaultConfigPath returns the default path for the configuration file
func DefaultConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "cdns", "config.yaml"), nil
}

// EnsureConfigFile ensures the config file exists, creating it with defaults if not
func EnsureConfigFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create directory
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Write default config
		if err := os.WriteFile(path, []byte(strings.TrimSpace(DefaultConfig)), 0644); err != nil {
			return fmt.Errorf("failed to write default config: %w", err)
		}
	}
	return nil
}

// Load loads configuration from multiple sources with priority:
// 1. Flags (highest priority)
// 2. Environment variables (prefix: APP_)
// 3. Config file (if exists)
// 4. Defaults (lowest priority)
func (l *Loader) Load(configFile string, flags *pflag.FlagSet) (*Config, error) {
	// Set defaults first
	if err := l.loadDefaults(); err != nil {
		return nil, fmt.Errorf("failed to load defaults: %w", err)
	}

	// Load from config file if provided
	if configFile != "" {
		if err := l.k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("failed to load config file %s: %w", configFile, err)
		}
	}

	// Load from environment variables with APP_ prefix
	if err := l.k.Load(env.Provider("APP_", ".", func(s string) string {
		// Convert APP_FOO_BAR to foo.bar
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "APP_")), "_", ".")
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Load from flags (highest priority)
	if flags != nil {
		if err := l.k.Load(posflag.Provider(flags, ".", l.k), nil); err != nil {
			return nil, fmt.Errorf("failed to load flags: %w", err)
		}
	}

	// Unmarshal into strongly-typed Config struct
	var cfg Config
	if err := l.k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set loaded from path
	cfg.LoadedFrom = configFile

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// loadDefaults sets default configuration values
func (l *Loader) loadDefaults() error {
	defaults := map[string]interface{}{
		"logger.level":  "warn",
		"logger.format": "text",
	}

	for k, v := range defaults {
		if err := l.k.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}
