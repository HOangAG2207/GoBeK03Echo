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

	testCases := []struct {
		name string

		inputServiceName string
		inputInstanceID  string

		setupMockRepo func(ctx context.Context) *mocks.Repository

		expectedResponse *model.HealthCheckResponse
		expectedError    error
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
			expectedResponse: &model.HealthCheckResponse{
				Message:     "OK",
				ServiceName: "test-service",
				InstanceID:  "instance-1",
			},
			expectedError: nil,
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
			expectedResponse: nil,
			expectedError:    ErrRedisConnection,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRepo := tc.setupMockRepo(ctx)
			svc := NewService(mockRepo, tc.inputServiceName, tc.inputInstanceID)
			response, err := svc.CheckHealth(ctx)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, response)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}
