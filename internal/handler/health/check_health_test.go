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
		expectedBody   map[string]interface{}
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
			expectedBody: map[string]any{
				"status":  "success",
				"message": "Health check passed",
			},
		},
		{
			name: "error - return 500",
			setupMock: func(ctx context.Context) *mocks.Service {
				s := mocks.NewService(t)
				s.On("CheckHealth", ctx).Return(nil, errors.New("redis error"))
				return s
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]any{
				"status":  "fail",
				"message": "Internal Server Error",
			},
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

			var body map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &body)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedBody["status"], body["status"])
			assert.Equal(t, tc.expectedBody["message"], body["message"])

			// check thêm field data khi success
			if rec.Code == http.StatusOK {
				assert.NotNil(t, body["data"])
			}

			// check thêm field error khi fail
			if rec.Code == http.StatusInternalServerError {
				assert.NotNil(t, body["error"])
			}
		})
	}
}
