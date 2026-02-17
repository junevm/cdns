package reset

import (
	"context"
	"log/slog"
	"testing"

	"github.com/junevm/cdns/internal/dns/models"
	"github.com/junevm/cdns/internal/features/status"
	"github.com/junevm/cdns/internal/ui"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDetector is a mock of reset.Detector
type MockDetector struct {
	mock.Mock
}

func (m *MockDetector) Detect() (models.Backend, error) {
	args := m.Called()
	return args.Get(0).(models.Backend), args.Error(1)
}

// MockReader is a mock of reset.Reader
type MockReader struct {
	mock.Mock
}

func (m *MockReader) ReadDNSConfig(ctx context.Context, backend models.Backend) (*status.StatusInfo, error) {
	args := m.Called(ctx, backend)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*status.StatusInfo), args.Error(1)
}

// MockDNSWriter is a mock of reset.DNSWriter
type MockDNSWriter struct {
	mock.Mock
}

func (m *MockDNSWriter) ResetToAutomatic(ctx context.Context, backend models.Backend, interfaces []string) error {
	args := m.Called(ctx, backend, interfaces)
	return args.Error(0)
}

func TestResetService_Reset(t *testing.T) {
	t.Run("successful reset", func(t *testing.T) {
		backend := models.BackendNetworkManager
		interfaces := []status.InterfaceStatus{
			{Name: "eth0"},
			{Name: "wlan0"},
		}
		statusInfo := &status.StatusInfo{
			Backend:    backend,
			Interfaces: interfaces,
		}

		mockDetector := new(MockDetector)
		mockDetector.On("Detect").Return(backend, nil)

		mockReader := new(MockReader)
		mockReader.On("ReadDNSConfig", mock.Anything, backend).Return(statusInfo, nil)

		mockWriter := new(MockDNSWriter)
		mockWriter.On("ResetToAutomatic", mock.Anything, backend, []string{"eth0", "wlan0"}).Return(nil)

		svc := &Service{
			detector: mockDetector,
			reader:   mockReader,
			writer:   mockWriter,
			logger:   slog.Default(),
			styles:   ui.NewStyles(),
		}

		err := svc.Reset(context.Background())
		assert.NoError(t, err)
		mockDetector.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockWriter.AssertExpectations(t)
	})

	t.Run("no active interfaces", func(t *testing.T) {
		backend := models.BackendNetworkManager
		statusInfo := &status.StatusInfo{
			Backend:    backend,
			Interfaces: []status.InterfaceStatus{},
		}

		mockDetector := new(MockDetector)
		mockDetector.On("Detect").Return(backend, nil)

		mockReader := new(MockReader)
		mockReader.On("ReadDNSConfig", mock.Anything, backend).Return(statusInfo, nil)

		// Writer should NOT be called
		mockWriter := new(MockDNSWriter)

		svc := &Service{
			detector: mockDetector,
			reader:   mockReader,
			writer:   mockWriter,
			logger:   slog.Default(),
			styles:   ui.NewStyles(),
		}

		err := svc.Reset(context.Background())
		assert.NoError(t, err)
		mockDetector.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockWriter.AssertExpectations(t)
	})
}
