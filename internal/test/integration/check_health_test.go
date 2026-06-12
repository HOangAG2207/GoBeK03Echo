package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlerPkg "github.com/HOangAG2207/GoBeK03Echo/internal/handler/health"
	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	serviceMock "github.com/HOangAG2207/GoBeK03Echo/internal/service/health/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIntegration_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupHTTP func(e *echo.Echo) *httptest.ResponseRecorder

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Normal case - success",
			setupHTTP: func(e *echo.Echo) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)
				rec := httptest.NewRecorder()

				e.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"status":"success",
				"info":"Health check passed",
				"message":"OK",
				"service_name":"test-service",
				"instance_id":"instance-1"
			}`,
		},
		{
			name: "Not found route",
			setupHTTP: func(e *echo.Echo) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, "/unknown", nil)
				rec := httptest.NewRecorder()

				e.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"Not Found"}`,
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

			// setup behavior theo từng test
			if tc.expectedStatus == http.StatusOK {
				mockService.On("CheckHealth", mock.Anything).
					Return(&model.HealthCheckResponse{
						Message:     "OK",
						ServiceName: "test-service",
						InstanceID:  "instance-1",
					}, nil)
			}

			// ===== register handler =====
			h := handlerPkg.NewHandler(mockService)
			e.GET("/v1/health-check", h.CheckHealth)

			// ===== call API =====
			rec := tc.setupHTTP(e)

			// ===== assert status =====
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// ===== assert body =====
			if tc.expectedStatus == http.StatusOK {
				assert.JSONEq(t, tc.expectedBody, rec.Body.String())
			}

			if tc.expectedStatus == http.StatusNotFound {
				assert.Equal(t, tc.expectedBody, strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}
