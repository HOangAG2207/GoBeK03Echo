package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/HOangAG2207/GoBeK03Echo/internal/service/health/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *mocks.Service

		expectedStatus int

		// chỉ giữ expectation cần thiết
		expectSuccess bool
		expectError   bool
	}{
		{
			name: "success - return 200",
			setupMock: func(ctx context.Context) *mocks.Service {
				s := mocks.NewService(t)
				s.On("CheckHealth", ctx).Return(&model.HealthCheckResponse{
					Message:     "OK",
					ServiceName: "test-service",
					InstanceID:  "instance-1",
				}, nil)
				return s
			},
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name: "error - return 500",
			setupMock: func(ctx context.Context) *mocks.Service {
				s := mocks.NewService(t)
				s.On("CheckHealth", ctx).Return(nil, errors.New("redis error"))
				return s
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			ctx := req.Context()
			mockService := tc.setupMock(ctx)

			h := NewHandler(mockService)

			err := h.CheckHealth(c)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]any
			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// ===== SUCCESS CASE =====
			if tc.expectSuccess {
				assert.Equal(t, "success", resp["status"])
				assert.Equal(t, "Health check passed", resp["info"])

				// data bị flatten
				assert.Equal(t, "OK", resp["message"])
				assert.Equal(t, "test-service", resp["service_name"])
				assert.Equal(t, "instance-1", resp["instance_id"])
			}

			// ===== ERROR CASE =====
			if tc.expectError {
				assert.Equal(t, "error", resp["status"])
				assert.Equal(t, "Internal Server Error", resp["info"])

				assert.Contains(t, resp, "error")
				assert.NotEmpty(t, resp["error"])
			}
		})
	}
}
