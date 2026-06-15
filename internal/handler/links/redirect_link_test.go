package links

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HOangAG2207/GoBeK03Echo/internal/service/links"
	"github.com/HOangAG2207/GoBeK03Echo/internal/service/links/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_RedirectLink(t *testing.T) {
	t.Parallel()

	type fields struct {
		setupRequest     func(ctx echo.Context)
		setupMockService func(ctx context.Context) *mocks.Service
	}
	type expected struct {
		Url    string
		Status int
	}
	testCases := []struct {
		name     string
		fields   fields
		expected expected
	}{
		{
			name: "normal case - success",
			fields: fields{
				setupRequest: func(ctx echo.Context) {
					req := httptest.NewRequest(
						http.MethodGet,
						"/v1/links/redirect/abc123",
						nil,
					)
					ctx.SetRequest(req)
					ctx.SetParamNames("code")
					ctx.SetParamValues("abc123")
				},
				setupMockService: func(ctx context.Context) *mocks.Service {
					serviceMock := mocks.NewService(t)
					serviceMock.On("GetLink", ctx, "abc123").
						Return("https://google.com", nil).
						Once()
					return serviceMock
				},
			},
			expected: expected{
				Status: http.StatusFound,
				Url:    "https://google.com",
			},
		},
		{
			name: "missing code parameter",
			fields: fields{
				setupRequest: func(ctx echo.Context) {
					req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect", nil)
					ctx.SetRequest(req)
					// No code parameter set
				},

				setupMockService: func(ctx context.Context) *mocks.Service {
					serviceMock := mocks.NewService(t)

					serviceMock.
						On("GetLink", ctx, "").
						Return("", links.ErrCodeNotFound).
						Once()

					return serviceMock
				},
			},
			expected: expected{
				Status: http.StatusNotFound,
				Url:    "",
			},
		},
		{
			name: "internal error should return 500",
			fields: fields{
				setupRequest: func(ctx echo.Context) {
					req := httptest.NewRequest(
						http.MethodGet,
						"/v1/links/redirect/abc123",
						nil,
					)

					ctx.SetRequest(req)
					ctx.SetParamNames("code")
					ctx.SetParamValues("abc123")
				},

				setupMockService: func(ctx context.Context) *mocks.Service {
					serviceMock := mocks.NewService(t)

					serviceMock.
						On("GetLink", ctx, "abc123").
						Return("", errors.New("redis down")).
						Once()

					return serviceMock
				},
			},
			expected: expected{
				Status: http.StatusInternalServerError,
				Url:    "",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := e.NewContext(req, rec)
			tc.fields.setupRequest(ctx)

			mockSvc := tc.fields.setupMockService(ctx.Request().Context())

			testHandler := NewHandler(mockSvc)

			testHandler.RedirectLink(ctx)

			assert.Equal(t, tc.expected.Status, rec.Code)
			assert.Equal(t, tc.expected.Url, rec.Header().Get("Location"))
		})
	}
}
