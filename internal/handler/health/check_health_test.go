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
	"github.com/HOangAG2207/GoBeK03Echo/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *mocks.Service

		expectedStatus int

		expectedSuccess *utils.SuccessResponse
		expectedError   *utils.ErrorResponse
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
			expectedSuccess: &utils.SuccessResponse{
				Status:  "success",
				Message: "Health check passed",
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

			if tc.expectedSuccess != nil {
				var resp utils.SuccessResponse
				err = json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedSuccess.Status, resp.Status)
				assert.Equal(t, tc.expectedSuccess.Message, resp.Message)

				// parse data field sâu hơn
				dataBytes, _ := json.Marshal(resp.Data)
				var data model.HealthCheckResponse
				err = json.Unmarshal(dataBytes, &data)
				assert.NoError(t, err)

				assert.Equal(t, "OK", data.Message)
				assert.Equal(t, "test-service", data.ServiceName)
				assert.Equal(t, "instance-1", data.InstanceID)
			}

			if tc.expectedError != nil {
				var resp utils.ErrorResponse
				err = json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedError.Status, resp.Status)
				assert.Equal(t, tc.expectedError.Message, resp.Message)

				assert.NotEmpty(t, resp.Error)
			}
		})
	}
}
