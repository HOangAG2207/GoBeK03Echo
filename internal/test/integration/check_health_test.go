package integration_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handlerPkg "github.com/HOangAG2207/GoBeK03Echo/internal/handler/health"
	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	serviceMock "github.com/HOangAG2207/GoBeK03Echo/internal/service/health/mocks"
	"github.com/HOangAG2207/GoBeK03Echo/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIntegration_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(s *serviceMock.Service)

		expectedStatus int

		expectedSuccess *utils.SuccessResponse
		expectedError   *utils.ErrorResponse
	}{
		{
			name: "success - full flow OK",
			setupMock: func(s *serviceMock.Service) {
				s.On("CheckHealth", mock.Anything).
					Return(&model.HealthCheckResponse{
						Message:     "OK",
						ServiceName: "test-service",
						InstanceID:  "instance-1",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedSuccess: &utils.SuccessResponse{
				Status:  "success",
				Message: "Health check passed",
			},
		},
		{
			name: "error - service fail",
			setupMock: func(s *serviceMock.Service) {
				s.On("CheckHealth", mock.Anything).
					Return(nil, errors.New("something went wrong"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError: &utils.ErrorResponse{
				Status:  "fail",
				Message: "Internal Server Error",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// ===== setup echo =====
			e := echo.New()

			// ===== mock service =====
			mockService := serviceMock.NewService(t)
			tc.setupMock(mockService)

			// ===== handler =====
			handler := handlerPkg.NewHandler(mockService)

			e.GET("/v1/health-check", handler.CheckHealth)

			// ===== request =====
			req := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			// ===== assert status =====
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// ===== success case =====
			if tc.expectedSuccess != nil {
				var resp utils.SuccessResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedSuccess.Status, resp.Status)
				assert.Equal(t, tc.expectedSuccess.Message, resp.Message)

				// ✅ strong typing cho data
				dataBytes, _ := json.Marshal(resp.Data)

				var data model.HealthCheckResponse
				err = json.Unmarshal(dataBytes, &data)
				assert.NoError(t, err)

				assert.Equal(t, "OK", data.Message)
				assert.Equal(t, "test-service", data.ServiceName)
				assert.Equal(t, "instance-1", data.InstanceID)
			}

			// ===== error case =====
			if tc.expectedError != nil {
				var resp utils.ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedError.Status, resp.Status)
				assert.Equal(t, tc.expectedError.Message, resp.Message)
				assert.NotEmpty(t, resp.Error)
			}
		})
	}
}
