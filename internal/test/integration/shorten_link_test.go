package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlerLinks "github.com/HOangAG2207/GoBeK03Echo/internal/handler/links"
	serviceMock "github.com/HOangAG2207/GoBeK03Echo/internal/service/links/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIntegration_ShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(s *serviceMock.Service)
		setupHTTP func(e *echo.Echo) *httptest.ResponseRecorder

		expectedStatus  int
		expectedCodeLen int
	}{
		{
			name: "success - shorten ok",
			setupMock: func(s *serviceMock.Service) {
				s.On("ShortenLink", mock.Anything, "http://example.com", 7, int64(3600)).
					Return("abc123xyz", nil)
			},
			setupHTTP: func(e *echo.Echo) *httptest.ResponseRecorder {
				body := `{"url":"http://example.com","exp":3600}`
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus:  http.StatusOK,
			expectedCodeLen: 9,
		},
		{
			name: "invalid payload",
			setupMock: func(s *serviceMock.Service) {
				// không gọi service vì bind fail
			},
			setupHTTP: func(e *echo.Echo) *httptest.ResponseRecorder {
				body := `{"url":"", "exp":-1}`
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus:  http.StatusBadRequest,
			expectedCodeLen: 0,
		},
		{
			name: "service error",
			setupMock: func(s *serviceMock.Service) {
				s.On("ShortenLink", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return("", fmt.Errorf("redis down"))
			},
			setupHTTP: func(e *echo.Echo) *httptest.ResponseRecorder {
				body := `{"url":"http://example.com","exp":3600}`
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedCodeLen: 0,
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
			if tc.setupMock != nil {
				tc.setupMock(mockService)
			}

			// ===== register handler =====
			h := handlerLinks.NewHandler(mockService)
			e.POST("/v1/links/shorten", h.ShortenLink)

			// ===== call API =====
			rec := tc.setupHTTP(e)

			// ===== assert status =====
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// ===== parse response =====
			var resp map[string]any
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// ===== SUCCESS =====
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, "success", resp["status"])

				code, ok := resp["code"].(string)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCodeLen, len(code))
			}

			// ===== ERROR =====
			if tc.expectedStatus != http.StatusOK {
				assert.Equal(t, "error", resp["status"])
				assert.Contains(t, resp, "error")
			}
		})
	}
}
