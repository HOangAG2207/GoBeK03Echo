package health

import (
	"context"
	"errors"
	"testing"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/HOangAG2207/GoBeK03Echo/internal/repository/health/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_CheckHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		inputServiceName string
		inputInstanceID  string

		setupMockRepo func(ctx context.Context) *mocks.Repository

		expectError bool
		checkResp   func(t *testing.T, resp *model.HealthCheckResponse)
	}{
		{
			name:             "check health return OK status",
			inputServiceName: "test-service",
			inputInstanceID:  "instance-1",
			setupMockRepo: func(ctx context.Context) *mocks.Repository {
				repo := mocks.NewRepository(t)
				repo.On("PingRedis", ctx).Return(nil)
				return repo
			},
			expectError: false,
			checkResp: func(t *testing.T, resp *model.HealthCheckResponse) {
				assert.Equal(t, "OK", resp.Message)
				assert.Equal(t, "test-service", resp.ServiceName)
				assert.Equal(t, "instance-1", resp.InstanceID)
			},
		},

		{
			name:             "check health return redis error",
			inputServiceName: "test-service",
			inputInstanceID:  "instance-1",
			setupMockRepo: func(ctx context.Context) *mocks.Repository {
				repo := mocks.NewRepository(t)
				repo.On("PingRedis", ctx).Return(errors.New("redis down"))
				return repo
			},
			expectError: true,
			checkResp: func(t *testing.T, resp *model.HealthCheckResponse) {
				assert.Nil(t, resp)
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			mockRepo := tc.setupMockRepo(ctx)

			svc := NewService(mockRepo, tc.inputServiceName, tc.inputInstanceID)

			resp, err := svc.CheckHealth(ctx)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)

			if tc.checkResp != nil {
				tc.checkResp(t, resp)
			}
		})
	}
}
