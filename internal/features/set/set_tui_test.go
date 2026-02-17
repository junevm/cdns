package set

import (
	"github.com/junevm/cdns/internal/config"
	"github.com/junevm/cdns/internal/dns/backend"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestService_InteractiveMode(t *testing.T) {
	s := NewService(&config.Config{}, slog.Default(), &backend.DefaultSystemOps{})

	t.Run("interactive set mode exists", func(t *testing.T) {
		assert.NotNil(t, s)
		// We will test RunInteractiveSet in more detail as we implement it
	})
}
