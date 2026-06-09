package links

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HOangAG2207/GoBeK03Echo/internal/service/links/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ShortenURL(t *testing.T) {
	t.Parallel()

	type fields struct {
		setupRequest     func(ctx echo.Context)
		setupMockService func(ctx context.Context) *mocks.Service
	}
	type expected struct {
		Status   int
		Response string
	}

	testCases := []struct {
		name     string
		fields   fields
		expected expected
	}{
		{
			name: "success",
			fields: fields{
				setupRequest: func(ctx echo.Context) {
					req := httptest.NewRequest(
						http.MethodPost,
						"/v1/links/shorten",
						strings.NewReader(`{"url":"https://google.com","exp":3600}`), // ✅ FIX exp
					)
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					ctx.SetRequest(req)
				},
				setupMockService: func(ctx context.Context) *mocks.Service {
					mockSvc := mocks.NewService(t)

					mockSvc.
						On("ShortenURL", mock.Anything, "https://google.com", 7, int64(3600)).
						Return("abc123", nil).
						Once()

					return mockSvc
				},
			},
			expected: expected{
				Status:   http.StatusOK,
				Response: `"code":"abc123"`, // ✅ chỉ check phần cần thiết
			},
		},
		{
			name: "bad_request_invalid_json",
			fields: fields{
				setupRequest: func(ctx echo.Context) {
					req := httptest.NewRequest(
						http.MethodPost,
						"/v1/links/shorten",
						strings.NewReader(`{"url":}`), // invalid JSON
					)
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					ctx.SetRequest(req)
				},
				setupMockService: func(ctx context.Context) *mocks.Service {
					return mocks.NewService(t) // không expect call
				},
			},
			expected: expected{
				Status:   http.StatusBadRequest,
				Response: `"invalid request body"`,
			},
		},
		{
			name: "bad_request_validation",
			fields: fields{
				setupRequest: func(ctx echo.Context) {
					req := httptest.NewRequest(
						http.MethodPost,
						"/v1/links/shorten",
						strings.NewReader(`{"url":"","exp":0}`),
					)
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					ctx.SetRequest(req)
				},
				setupMockService: func(ctx context.Context) *mocks.Service {
					return mocks.NewService(t) // không expect call
				},
			},
			expected: expected{
				Status:   http.StatusBadRequest,
				Response: `"validation error"`,
			},
		},
		{
			name: "internal_server_error",
			fields: fields{
				setupRequest: func(ctx echo.Context) {
					req := httptest.NewRequest(
						http.MethodPost,
						"/v1/links/shorten",
						strings.NewReader(`{"url":"https://google.com","exp":3600}`),
					)
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

					ctx.SetRequest(req)
				},
				setupMockService: func(ctx context.Context) *mocks.Service {
					mockSvc := mocks.NewService(t)

					mockSvc.
						On("ShortenURL", mock.Anything, "https://google.com", 7, int64(3600)).
						Return("", assert.AnError).
						Once()

					return mockSvc
				},
			},
			expected: expected{
				Status:   http.StatusInternalServerError,
				Response: `"Internal Server Error"`, // ✅ FIX casing
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", nil)

			ctx := e.NewContext(req, rec)

			// setup request
			tc.fields.setupRequest(ctx)

			// setup mock
			mockSvc := tc.fields.setupMockService(ctx.Request().Context())

			handler := NewHandler(mockSvc)

			// call handler
			err := handler.ShortenURL(ctx)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expected.Status, rec.Code)

			body := rec.Body.String()
			assert.Contains(t, body, tc.expected.Response)

			mockSvc.AssertExpectations(t)
		})
	}
}
