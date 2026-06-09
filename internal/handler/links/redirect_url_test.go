package links

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HOangAG2207/GoBeK03Echo/internal/service/links"
	"github.com/HOangAG2207/GoBeK03Echo/internal/service/links/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_RedirectUrl(t *testing.T) {
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
					serviceMock.On("GetURL", ctx, "abc123").
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
						On("GetURL", ctx, "").
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

			testHandler.RedirectURL(ctx)

			assert.Equal(t, tc.expected.Status, rec.Code)
			assert.Equal(t, tc.expected.Url, rec.Header().Get("Location"))
		})
	}
}
