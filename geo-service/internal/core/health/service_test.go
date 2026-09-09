package health_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/telemetry-platform/geo-service/internal/core/health"
	"github.com/telemetry-platform/geo-service/internal/infrastructure/pkg/logger"
)

func TestService_GetHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "works correctly"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := health.NewService(logger.NewLogger())

			result, err := svc.GetHealth(context.Background())

			assert.NoError(t, err)
			assert.Equal(t, "ok", result.Status)
			assert.Equal(t, "geo-service", result.Service)
		})
	}
}
