package list

import (
	"log/slog"
	"sort"
	"testing"

	"cli/internal/config"

	"github.com/stretchr/testify/assert"
)

// ListPresets returns a list of DNS presets combined from built-ins and custom config
func TestListPresets(t *testing.T) {
	// Setup custom config
	cfg := &config.Config{
		DNS: config.DNSConfig{
			CustomPresets: map[string][]string{
				"my-custom": {"1.2.3.4", "5.6.7.8"},
			},
		},
	}
	logger := slog.Default()

	// Create service
	svc := NewService(cfg, logger)

	// Invoke ListPresets
	presets := svc.ListPresets()

	// Assertions
	// 1. Should contain built-in presets (e.g., Google or Cloudflare)
	foundCloudflare := false
	for _, p := range presets {
		if p.Name == "Cloudflare" {
			foundCloudflare = true
			assert.Equal(t, "Built-in", p.Type)
			break
		}
	}
	assert.True(t, foundCloudflare, "Should find built-in preset Cloudflare")

	// 2. Should contain custom preset
	foundCustom := false
	for _, p := range presets {
		if p.Name == "my-custom" {
			foundCustom = true
			assert.Equal(t, "Custom", p.Type)
			assert.Equal(t, "1.2.3.4, 5.6.7.8", p.Servers)
			break
		}
	}
	assert.True(t, foundCustom, "Should find custom preset my-custom")
}

// Helper to check sorting if applicable, though map iteration order is random
func TestListPresets_Sorting(t *testing.T) {
	cfg := &config.Config{}
	svc := NewService(cfg, slog.Default())
	presets := svc.ListPresets()

	// Basic check that we get something back
	assert.NotEmpty(t, presets)

	// Check if results are sorted by name for display consistency?
	// The implementation should probably sort them.
	// Let's assume we want them sorted.
	isSorted := sort.SliceIsSorted(presets, func(i, j int) bool {
		return presets[i].Name < presets[j].Name
	})
	// This might fail initially if implementation doesn't sort, which is fine for TDD loop.
	assert.True(t, isSorted, "Presets should be sorted by name")
}
