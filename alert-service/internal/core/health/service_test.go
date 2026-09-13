package health_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/telemetry-platform/alert-service/internal/core/health"
	repomocks "github.com/telemetry-platform/alert-service/internal/core/ports/repositories/mocks"
	"github.com/telemetry-platform/alert-service/internal/infrastructure/pkg/logger"
)

func TestService_GetHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func(*repomocks.Databaser)
		expectedStatus string
	}{
		{
			name: "returns ok when database is reachable",
			setup: func(m *repomocks.Databaser) {
				m.On("PingContext", mock.Anything).Return(nil)
			},
			expectedStatus: "ok",
		},
		{
			name: "returns degraded when database is unreachable",
			setup: func(m *repomocks.Databaser) {
				m.On("PingContext", mock.Anything).Return(assert.AnError)
			},
			expectedStatus: "degraded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db := repomocks.NewDatabaser(t)
			tt.setup(db)
			svc := health.NewService(db, logger.NewLogger())

			result, err := svc.GetHealth(context.Background())

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, result.Status)
			assert.Equal(t, "alert-service", result.Service)
		})
	}
}
